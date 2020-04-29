// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package randdata_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/tunabay/go-infounit"
	"github.com/tunabay/go-randdata"
)

//
func TestReader_new(t *testing.T) {
	t.Parallel()

	buf := make([]byte, 32)
	r := randdata.New(randdata.Binary, 123, 64*1000+13)
	v := randdata.NewVerifier(randdata.Binary, 123, 64*1000+13 /*-8*/)
	for i := 0; ; i++ {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Error(err)
			break
		}
		if n == 0 {
			t.Errorf("zero")
			continue
		}
		// t.Logf("BIN: %4d: %5d: %x\n", i, n, buf[:n])

		wn, err := v.Write(buf[:n])
		if err != nil {
			t.Errorf("verify failed: %s", err)
			continue
		}
		_ = wn
		// t.Logf("%d B verified\n", wn)
	}
}

//
func TestReader_2(t *testing.T) {
	t.Parallel()

	tct := []randdata.Type{
		randdata.Zero,
		randdata.Binary,
		randdata.Text,
		randdata.Type(127),
	}
	tcs := []infounit.ByteCount{
		0, 1, 2, 7, 256, 1023, 8192,
		infounit.Megabyte * 1,
	}
	for _, tt := range tct {
		tt := tt
		for i, ts := range tcs {
			ts := ts
			t.Run(fmt.Sprintf("%s_%s", tt, ts), func(t *testing.T) {
				t.Parallel()
				r := randdata.New(tt, int64(1234+i), ts)
				v := randdata.NewVerifier(tt, int64(1234+i), ts)
				if _, err := r.WriteTo(v); err != nil {
					t.Errorf("%v %v: verify failed: %s", tt, ts, err)
				}
				if err := v.Close(); err != nil {
					t.Errorf("%v %v: verify failed at close: %s", tt, ts, err)
				}
			})
		}
	}
}

//
func TestReader_newAsFile(t *testing.T) {
	t.Parallel()

	tct := []randdata.Type{randdata.Binary, randdata.Text}
	tcs := []infounit.ByteCount{
		0, 1, 2, 7, 256, 1023, 8192,
		infounit.Megabyte * 1,
	}
	for _, tt := range tct {
		tt := tt
		for i, ts := range tcs {
			ts := ts
			t.Run(fmt.Sprintf("%s_%s", tt, ts), func(t *testing.T) {
				t.Parallel()
				fn := fmt.Sprintf("/tmp/tr-test-%s-%s.test.tmp", tt, ts)
				v, err := randdata.NewAsFile(tt, int64(1234+i), ts, fn)
				if err != nil {
					t.Errorf("%v %v: NewAsFile() failed: %s", tt, ts, err)
					return
				}
				defer os.Remove(fn)
				f, err := os.Open(fn)
				if err != nil {
					t.Errorf("%v %v: Open %s failed: %s", tt, ts, fn, err)
					return
				}
				defer f.Close()
				n, err := v.ReadFrom(f)
				if err != nil {
					t.Errorf("%v %v: verify failed: %s", tt, ts, err)
					return
				}
				if vs := infounit.ByteCount(n); vs != ts {
					t.Errorf("%v %v: unexpected verified len: want: %d, got: %d", tt, ts, ts, vs)
					return
				}
				if err := v.Close(); err != nil {
					t.Errorf("%v %v: verify failed at close: %s", tt, ts, err)
					return
				}
			})
		}
	}
}

//
func TestReader_newAsTempFile(t *testing.T) {
	t.Parallel()

	tct := []randdata.Type{randdata.Binary, randdata.Text}
	tcs := []infounit.ByteCount{
		0, 1, 2, 7, 256, 1023, 8192,
		infounit.Megabyte * 1,
	}
	for _, tt := range tct {
		tt := tt
		for i, ts := range tcs {
			ts := ts
			t.Run(fmt.Sprintf("%s_%s", tt, ts), func(t *testing.T) {
				t.Parallel()
				fn, v, err := randdata.NewAsTempFile(tt, int64(1234+i), ts)
				if err != nil {
					t.Errorf("%v %v: NewAsTempFile() failed: %s", tt, ts, err)
					return
				}
				defer os.Remove(fn)
				n, err := v.ReadFromFile(fn)
				if err != nil {
					t.Errorf("%v %v: verify failed: %s", tt, ts, err)
					return
				}
				if vs := infounit.ByteCount(n); vs != ts {
					t.Errorf("%v %v: unexpected verified len: want: %d, got: %d", tt, ts, ts, vs)
					return
				}
				if err := v.Close(); err != nil {
					t.Errorf("%v %v: verify failed at close: %s", tt, ts, err)
					return
				}
			})
		}
	}
}
