# go-randdata

[![GitHub](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/tunabay/go-randdata/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/tunabay/go-randdata?status.svg)](https://godoc.org/github.com/tunabay/go-randdata)
[![Go Report Card](https://goreportcard.com/badge/github.com/tunabay/go-randdata)](https://goreportcard.com/report/github.com/tunabay/go-randdata)
[![codecov](https://codecov.io/gh/tunabay/go-randdata/branch/master/graph/badge.svg)](https://codecov.io/gh/tunabay/go-randdata)

go-randdata is a Go package providing a mechanism for unit testing to generate
and verify reproducible pseudo-random byte sequences.

Reader is the pseudo-random byte sequence generator. It implements the io.Reader
interface and can be Read the generated byte sequence. Verifier is the Reader
companion object that implements the io.Writer interface. It verifies that the
data written is exactly the same as the byte sequence generated by the Reader.
```
import (
	"fmt"
	"io"
	"github.com/tunabay/go-randdata"
)

func main() {
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
}
```
[Run in Go Playground](https://play.golang.org/p/SgEvIZ0cNjw)

The Reader also generates "jitter" to reading operation. In the above example,
calling Read method with the 256 bytes buffer returns randomly shorter written
length. While the Read method of the io.Reader interface can return shorter
length than passed buffer, program should be able to handle that.

## Documentation

- http://godoc.org/github.com/tunabay/go-randdata

## See also

- https://golang.org/pkg/testing/

## License

go-randdata is available under the MIT license. See the [LICENSE](LICENSE) file for more information.
