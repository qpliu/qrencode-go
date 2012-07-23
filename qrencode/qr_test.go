package qrencode

import (
	"image/png"
	"os"
	"testing"
)

func TestContentBits(t *testing.T) {
	bits, version, err := stringContentBits("HELLO WORLD", ECLevelQ)
	if err != nil {
		t.Error(err.Error())
	}
	if version != versionNumber(1) {
		t.Error("version", version, " != 1")
	}
	if bits.String() != "00100000010110110000101101111000110100010111001011011100010011010100001101000000111011000001000111101100" {
		t.Error("bits", bits.String(), " != 00100000010110110000101101111000110100010111001011011100010011010100001101000000111011000001000111101100")
	}
	bits = interleaveWithECBytes(bits, version, ECLevelQ)
	if bits.String() != "0010000001011011000010110111100011010001011100101101110001001101010000110100000011101100000100011110110010101000010010000001011001010010110110010011011010011100000000000010111000001111101101000111101000010000" {
		t.Error("bits", bits.String(), " != 0010000001011011000010110111100011010001011100101101110001001101010000110100000011101100000100011110110010101000010010000001011001010010110110010011011010011100000000000010111000001111101101000111101000010000")
	}
}

func TestGenerateECBytes(t *testing.T) {
	block := blockPair{
		dataBytes: []int{32, 91, 11, 120, 209, 114, 220, 77, 67, 64, 236, 17, 236},
		ecBytes:   make([]int, 13),
	}
	generateECBytes(&block)
	for i, b := range []int{32, 91, 11, 120, 209, 114, 220, 77, 67, 64, 236, 17, 236} {
		if block.dataBytes[i] != b {
			t.Error("dataBytes", i, block.dataBytes[i], b)
		}
	}
	for i, b := range []int{168, 72, 22, 82, 217, 54, 156, 0, 46, 15, 180, 122, 16} {
		if block.ecBytes[i] != b {
			t.Error("ecBytes", i, block.ecBytes[i], b)
		}
	}
}

func BenchmarkEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Encode("Testing one two three four five six seven eight nine ten eleven twelve thirteen", ECLevelQ)
	}
}

func ExampleEncode() {
	grid, err := Encode("Testing one two three four five six seven eight nine ten eleven twelve thirteen fourteen fifteen sixteen seventeen eighteen nineteen twenty.", ECLevelQ)
	if err != nil {
		return
	}
	f, err := os.Create("/tmp/qr.png")
	if err != nil {
		return
	}
	defer f.Close()
	png.Encode(f, grid.Image(8))
}
