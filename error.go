// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package randdata

import (
	"fmt"

	"github.com/tunabay/go-infounit"
)

// TrailingExtraBytesError is the error thrown when extra bytes are written to
// the Verifier beyond the expected data end.
type TrailingExtraBytesError struct {
	ExpectedLen   infounit.ByteCount
	WrittenLen    infounit.ByteCount
	AcceptedBytes []byte
	ExtraBytes    []byte
}

// Error returns the string representation of TrailingExtraBytesError.
func (err TrailingExtraBytesError) Error() string {
	return fmt.Sprintf("trailing extra bytes: %x", err.ExtraBytes)
}

// NotEnoughBytesError is the error thrown when the Verifier is closed before
// the expected size of data is written.
type NotEnoughBytesError struct {
	ExpectedLen infounit.ByteCount
	WrittenLen  infounit.ByteCount
}

// Error returns the string representation of NotEnoughBytesError.
func (err NotEnoughBytesError) Error() string {
	return fmt.Sprintf("not enough bytes: %v short", err.ExpectedLen-err.WrittenLen)
}

// UnexpectedBytesError is the error thrown when an unexpected wrong byte
// sequence is written to the Verifier.
type UnexpectedBytesError struct {
	Pos           infounit.ByteCount // the first byte is 0
	ExpectedBytes []byte
	WrittenBytes  []byte
}

// Error returns the string representation of UnexpectedBytesError.
func (err UnexpectedBytesError) Error() string {
	return fmt.Sprintf("unexpected byte at %d: want: %x, got: %x", err.Pos, err.ExpectedBytes, err.WrittenBytes)
}
