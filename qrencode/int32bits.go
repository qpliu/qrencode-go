package qrencode

type uint32BitVector struct {
	bitIndex byte
	bits     []uint32
}

func (v *uint32BitVector) Length() int {
	if v.bitIndex == 0 {
		return len(v.bits) * 32
	}
	return (len(v.bits)-1)*32 + int(v.bitIndex)
}

func (v *uint32BitVector) Get(i int) bool {
	return (v.bits[i/32]>>uint(i%32))&1 == 1
}

func (v *uint32BitVector) AppendBit(b bool) {
	if v.bitIndex == 0 {
		v.bits = append(v.bits, 0)
	}
	if b {
		v.bits[len(v.bits)-1] |= 1 << v.bitIndex
	} else {
		v.bits[len(v.bits)-1] &= 0xffffffff ^ (1 << v.bitIndex)
	}
	v.bitIndex++
	v.bitIndex %= 32
}

func (v *uint32BitVector) Append(b, count int) {
	for i := uint(count); i > 0; i-- {
		v.AppendBit((b>>(i-1))&1 == 1)
	}
}

func (v *uint32BitVector) AppendBits(b uint32BitVector) {
	if v.bitIndex == 0 {
		v.bitIndex = b.bitIndex
		v.bits = append(v.bits, b.bits...)
	} else {
		for i, l := 0, b.Length(); i < l; i++ {
			v.AppendBit(b.Get(i))
		}
	}
}

type uint32BitGrid struct {
	width, height int
	bits          []uint32
}

func newUint32BitGrid(width, height int) uint32BitGrid {
	return uint32BitGrid{width, height, make([]uint32, (2*width*height+31)/32)}
}

func (g *uint32BitGrid) Width() int {
	return g.width
}

func (g *uint32BitGrid) Height() int {
	return g.height
}

func (g *uint32BitGrid) Empty(x, y int) bool {
	i := 2 * (x + y*g.width)
	return (g.bits[i/32]>>uint(i%32))&1 == 0
}

func (g *uint32BitGrid) Get(x, y int) bool {
	i := 2*(x+y*g.width) + 1
	return (g.bits[i/32]>>uint(i%32))&1 == 0
}

func (g *uint32BitGrid) Set(x, y int, v bool) {
	i := 2 * (x + y*g.width)
	if v {
		g.bits[i/32] |= 3 << uint(i%32)
	} else {
		g.bits[i/32] |= 1 << uint(i%32)
		g.bits[i/32] &= 0xffffffff ^ (1 << uint(i%32+1))
	}
}

func (g *uint32BitGrid) Clear() {
	for i, _ := range g.bits {
		g.bits[i] = 0
	}
}
