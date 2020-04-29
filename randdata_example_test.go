// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package randdata_test

import (
	"fmt"
	"io"

	"github.com/tunabay/go-randdata"
)

//
func Example() {

	// 5 MB pseudo-random byte sequence, using random seed 123
	r := randdata.New(randdata.Binary, 123, 5000000)

	// paired verifier
	v := r.NewVerifier()

	// read and veriry data
	buf := make([]byte, 256)
	for {
		n, err := r.Read(buf)
		if 0 < n {
			if _, err := v.Write(buf[:n]); err != nil {
				fmt.Println(err)
				break
			}
		}
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
	}

	// verify that written data is enough
	if err := v.Close(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Read:", r.TotalRead())

	// Output:
	// Read: 5.0 MB
}

//
func Example_shortData() {

	// 3 MB pseudo-random byte sequence, using random seed 777
	r := randdata.New(randdata.Binary, 777, 3000000)

	// verifier expecting 10 bytes extra
	v := randdata.NewVerifier(randdata.Binary, 777, 3000000+10)

	// read and veriry data
	buf := make([]byte, 256)
	for {
		n, err := r.Read(buf)
		if 0 < n {
			if _, err := v.Write(buf[:n]); err != nil {
				fmt.Println(err)
				break
			}
		}
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
	}

	// verify that written data is enough
	if err := v.Close(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Read:", r.TotalRead())

	// Output:
	// not enough bytes: 10 B short
	// Read: 3.0 MB
}
