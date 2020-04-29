// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package randdata

import (
	"math/rand"

	"github.com/tunabay/go-infounit"
)

//
type zeroGenerator struct {
	blockSize infounit.ByteCount
}

//
func newZeroGenerator() *zeroGenerator {
	return &zeroGenerator{
		blockSize: 4096,
	}
}

//
func (g *zeroGenerator) Gen(r *rand.Rand, pos, rem infounit.ByteCount) ([]byte, error) {
	s := int(g.blockSize)
	if rem < g.blockSize {
		s = int(rem)
	}
	return make([]byte, s), nil
}
