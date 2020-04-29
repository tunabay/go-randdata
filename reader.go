// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package randdata

import (
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"sync"

	"github.com/tunabay/go-infounit"
)

// Reader is a reproducible and verifiable pseudo-random byte sequence
// generator. Reader with the same parameters, type, seed, and size, always
// produces the same byte sequence. The byte sequence read from Reader is
// verified by writing it to Verifier. Reader implements io.Reader,
// io.ByteReader, and io.WriterTo interfaces.
type Reader struct {
	dType  Type               // data type
	seed   int64              // random seed for data generation and length fluctuation
	size   infounit.ByteCount // data size
	read   infounit.ByteCount // number of bytes already read
	lenMin int                // minimum number of bytes read in one Read()
	lenMax int                // maximum number of bytes read in one Read()
	lenRnd *rand.Rand         // random generater for data
	datRnd *rand.Rand         // random generator for length
	gen    Generator          // data generator
	buf    []byte             // data currently being read
	mu     sync.Mutex
}

// NewAsFile generates a pseudo-random byte sequence into the specified file.
// The contents of the file can be verified by Verifier. For convenience, it
// returns the Verifier with identical parameters. It is the caller's
// responsibility to remove the file when no longer needed.
func NewAsFile(dType Type, seed int64, size infounit.ByteCount, path string) (*Verifier, error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	r := New(dType, seed, size)
	if _, err := r.WriteTo(f); err != nil {
		_ = f.Close()
		_ = os.Remove(f.Name())
		return nil, err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(f.Name())
		return nil, err
	}
	return r.NewVerifier(), nil
}

// NewAsTempFile generates a pseudo-random byte sequence into a temporary file,
// and returns the file name and the Verifier with identical parameters. The
// contents of the temporary file can be verified by Verifier. It is the
// caller's responsibility to remove the file when no longer needed.
func NewAsTempFile(dType Type, seed int64, size infounit.ByteCount) (string, *Verifier, error) {
	f, err := ioutil.TempFile("", "randdata*.tmp")
	if err != nil {
		return "", nil, err
	}
	r := New(dType, seed, size)
	if _, err := r.WriteTo(f); err != nil {
		_ = f.Close()
		_ = os.Remove(f.Name())
		return "", nil, err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(f.Name())
		return "", nil, err
	}
	return f.Name(), r.NewVerifier(), nil
}

// New creates and returns a new Reader instance with the specified data type,
// random seed, and size.
func New(dType Type, seed int64, size infounit.ByteCount) *Reader {
	r := &Reader{
		dType:  dType,
		seed:   seed,
		size:   size,
		lenMin: 1,
		lenMax: math.MaxInt32,
		lenRnd: rand.New(rand.NewSource(seed)),
		datRnd: rand.New(rand.NewSource(seed)),
	}
	switch dType {
	case Zero:
		r.gen = newZeroGenerator()
	case Binary:
		r.gen = newBinaryGenerator()
	case Text:
		r.gen = newTextGenerator()
	default:
		r.gen = newBinaryGenerator() // treated as Binary type
	}
	return r
}

// NewWithGenerator creates a new Reader instance with a custom Generator.
func NewWithGenerator(gen Generator, seed int64, size infounit.ByteCount) *Reader {
	return &Reader{
		dType:  Custom,
		seed:   seed,
		size:   size,
		lenMin: 1,
		lenMax: math.MaxInt32,
		lenRnd: rand.New(rand.NewSource(seed)),
		datRnd: rand.New(rand.NewSource(seed)),
		gen:    gen,
	}
}

// TODO: func (r *Reader) Signature() string

// TotalRead returns the total number of bytes already read.
func (r *Reader) TotalRead() infounit.ByteCount {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.read
}

// IsEOF returns true if the end of the pseudo-random byte sequence is reached.
func (r *Reader) IsEOF() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.size <= r.read
}

// NewVerifier creates a Verifier that can be used to verify byte sequence
// generated from this Reader.
func (r *Reader) NewVerifier() *Verifier {
	if r.dType == Custom {
		return NewVerifierWithGenerator(r.gen, r.seed, r.size)
	}
	return NewVerifier(r.dType, r.seed, r.size)
}

// WriteTo writes generated pseudo-random byte sequence to w. This implements
// the WriterTo interface in the package io.
func (r *Reader) WriteTo(w io.Writer) (int64, error) {
	const bufSize = 8192
	var wLen int64
	buf := make([]byte, bufSize)
	for {
		rn, err := r.Read(buf)
		if err != nil {
			if err != io.EOF {
				return wLen, err // Read never return this error
			}
			break
		}
		wn, err := w.Write(buf[:rn])
		wLen += int64(wn)
		if err != nil {
			return wLen, err
		}
	}
	return wLen, nil
}

// ReadByte reads and returns the next byte from the generated pseudo-random
// byte sequence. This implements the ByteReader interface in the package io.
func (r *Reader) ReadByte() (byte, error) {
	buf := make([]byte, 1)
	if _, err := r.Read(buf); err != nil {
		return 0, err
	}
	return buf[0], nil
}

// prepareNext generates the following pseudo-random byte sequence. This is
// called when the buffer is empty.
func (r *Reader) prepareNext() error {
	buf, err := r.gen.Gen(r.datRnd, r.read, r.size-r.read)
	if err != nil {
		return err
	}
	r.buf = buf
	return nil
}

// Read reads up to len(p) bytes of generated pseudo-random byte sequence into
// p. It returns (0, io.EOF) at end of the byte sequence. According to the
// io.Reader interface specification, Read may return only a length less than
// len(p). To test this behavior, this Read is implemented to intentionally
// return short data using pseudo-random. Programs that do not handle this
// behavior properly should wrap Reader with bufio.Reader.
func (r *Reader) Read(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.size <= r.read {
		return 0, io.EOF
	}

	readLen := len(p)
	if r.lenMin < len(p) && r.lenMin < r.lenMax { // activate jitter
		max := r.lenMax
		if len(p) < r.lenMax {
			max = len(p)
		}
		readLen = r.lenMin + r.lenRnd.Intn(max-r.lenMin+1)
	}

	if r.size < r.read+infounit.ByteCount(readLen) {
		readLen = int(r.size - r.read)
	}

	remLen := readLen
	curDst := p
	for {
		if len(r.buf) == 0 {
			if err := r.prepareNext(); err != nil {
				return readLen - remLen, err
			}
		}
		n := len(r.buf)
		if remLen <= n {
			copy(curDst, r.buf[:remLen])
			r.read += infounit.ByteCount(remLen)
			r.buf = r.buf[remLen:]
			break
		}
		copy(curDst, r.buf)
		curDst = curDst[n:]
		r.read += infounit.ByteCount(n)
		remLen -= n
		r.buf = nil
	}
	return readLen, nil
}
