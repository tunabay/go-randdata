// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package randdata

import (
	"math"
)

// Type represents the type of the random data to be generated.
type Type uint8

const (
	Zero     Type = iota          // zero filled data (not random)
	Binary                        // binary data
	Text                          // text data
	UTF8Text                      // unicode text data (not implemented yet)
	Custom   Type = math.MaxUint8 // custon generator
)

// String returns the string representation of Type value.
func (t Type) String() string {
	switch t {
	case Zero:
		return "zero"
	case Binary:
		return "binary"
	case Text:
		return "text"
	case UTF8Text:
		return "utf8-text"
	case Custom:
		return "custom"
	}
	return "unknown"
}
