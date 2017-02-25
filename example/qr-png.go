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
