// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package randdata

import (
	"math/rand"

	"github.com/tunabay/go-infounit"
)

// Generator is the interface implemented by mechanisms that generate various
// random data. Gen() returns a pseudo-random byte sequence of arbitrary length
// using the passed pseudo-random generator r. The parameters pos and rem
// represent the position from the beginning of the data and the number of bytes
// remaining up to the end of the data, respectively. Gen() must return the same
// byte sequence for the same parameters to generate reproducible data.
// Generator must not hold any state. The byte sequence returned by Gen() should
// not depend on anything other than parameters.
type Generator interface {
	Gen(r *rand.Rand, pos, rem infounit.ByteCount) ([]byte, error)
}
