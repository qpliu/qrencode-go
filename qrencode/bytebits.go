package qrencode

type byteBitVector struct {
	bitIndex byte
	bits     []byte
}

func (v *byteBitVector) Length() int {
	if v.bitIndex == 0 {
		return len(v.bits) * 8
	}
	return (len(v.bits)-1)*8 + int(v.bitIndex)
}

func (v *byteBitVector) Get(i int) bool {
	return (v.bits[i/8]>>uint(i%8))&1 == 1
}

func (v *byteBitVector) AppendBit(b bool) {
	if v.bitIndex == 0 {
		v.bits = append(v.bits, 0)
	}
	if b {
		v.bits[len(v.bits)-1] |= 1 << v.bitIndex
	} else {
		v.bits[len(v.bits)-1] &= 255 ^ (1 << v.bitIndex)
	}
	v.bitIndex++
	v.bitIndex %= 8
}

func (v *byteBitVector) Append(b, count int) {
	for i := uint(count); i > 0; i-- {
		v.AppendBit((b>>(i-1))&1 == 1)
	}
}

func (v *byteBitVector) AppendBits(b byteBitVector) {
	if v.bitIndex == 0 {
		v.bitIndex = b.bitIndex
		v.bits = append(v.bits, b.bits...)
	} else {
		for i, l := 0, b.Length(); i < l; i++ {
			v.AppendBit(b.Get(i))
		}
	}
}

type byteBitGrid struct {
	width, height int
	bits          []byte
}

func newByteBitGrid(width, height int) byteBitGrid {
	return byteBitGrid{width, height, make([]byte, (2*width*height+7)/8)}
}

func (g *byteBitGrid) Width() int {
	return g.width
}

func (g *byteBitGrid) Height() int {
	return g.height
}

func (g *byteBitGrid) Empty(x, y int) bool {
	i := 2 * (x + y*g.width)
	return (g.bits[i/8]>>uint(i%8))&1 == 0
}

func (g *byteBitGrid) Get(x, y int) bool {
	i := 2*(x+y*g.width) + 1
	return (g.bits[i/8]>>uint(i%8))&1 == 0
}

func (g *byteBitGrid) Set(x, y int, v bool) {
	i := 2 * (x + y*g.width)
	if v {
		g.bits[i/8] |= 3 << uint(i%8)
	} else {
		g.bits[i/8] |= 1 << uint(i%8)
		g.bits[i/8] &= 255 ^ (1 << uint(i%8+1))
	}
}

func (g *byteBitGrid) Clear() {
	for i, _ := range g.bits {
		g.bits[i] = 0
	}
}
