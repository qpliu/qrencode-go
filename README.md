QR encoder in Go based on the ZXing encoder (http://code.google.com/p/zxing/).

[![GoDoc](https://godoc.org/github.com/qpliu/qrencode-go/qrencode?status.svg)](https://godoc.org/github.com/qpliu/qrencode-go/qrencode)
[![Build Status](https://travis-ci.org/qpliu/qrencode-go.svg?branch=master)](https://travis-ci.org/qpliu/qrencode-go)

I was surprised that I couldn't find a QR encoder in Go, especially
since the example at http://golang.org/doc/effective_go.html#web_server
is a QR code generator, though the QR encoding is done by an external
Google service in the example.

# Example

```go
package main

import (
	"bytes"
	"image/png"
	"os"

	"github.com/qpliu/qrencode-go/qrencode"
)

func main() {
	var buf bytes.Buffer
	for i, arg := range os.Args {
		if i > 1 {
			if err := buf.WriteByte(' '); err != nil {
				panic(err)
			}
		}
		if i > 0 {
			if _, err := buf.WriteString(arg); err != nil {
				panic(err)
			}
		}
	}
	grid, err := qrencode.Encode(buf.String(), qrencode.ECLevelQ)
	if err != nil {
		panic(err)
	}
	png.Encode(os.Stdout, grid.Image(8))
}
```
