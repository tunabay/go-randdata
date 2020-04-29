// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package randdata

import (
	"math/rand"

	"github.com/tunabay/go-infounit"
)

//
type binaryGenerator struct {
	blockSize infounit.ByteCount
}

//
func newBinaryGenerator() *binaryGenerator {
	return &binaryGenerator{
		blockSize: 256,
	}
}

//
func (g *binaryGenerator) Gen(r *rand.Rand, pos, rem infounit.ByteCount) ([]byte, error) {
	s := int(g.blockSize)
	if rem < g.blockSize {
		s = int(rem)
	}
	buf := make([]byte, s)
	_, _ = r.Read(buf) // always returns len(buf), nil
	return buf, nil
}
