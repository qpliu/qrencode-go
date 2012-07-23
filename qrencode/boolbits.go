package qrencode

type boolBitVector struct {
	bits []bool
}

func (v *boolBitVector) Length() int {
	return len(v.bits)
}

func (v *boolBitVector) Get(i int) bool {
	return v.bits[i]
}

func (v *boolBitVector) AppendBit(b bool) {
	v.bits = append(v.bits, b)
}

func (v *boolBitVector) Append(b, count int) {
	for i := uint(count); i > 0; i-- {
		v.AppendBit((b>>(i-1))&1 == 1)
	}
}

func (v *boolBitVector) AppendBits(b boolBitVector) {
	v.bits = append(v.bits, b.bits...)
}

type boolBitGrid struct {
	width, height int
	bits          []bool
}

func newBoolBitGrid(width, height int) boolBitGrid {
	return boolBitGrid{width, height, make([]bool, 2*width*height)}
}

func (g *boolBitGrid) Width() int {
	return g.width
}

func (g *boolBitGrid) Height() int {
	return g.height
}

func (g *boolBitGrid) Empty(x, y int) bool {
	return !g.bits[2*(x+y*g.width)]
}

func (g *boolBitGrid) Get(x, y int) bool {
	return g.bits[2*(x+y*g.width)+1]
}

func (g *boolBitGrid) Set(x, y int, v bool) {
	g.bits[2*(x+y*g.width)] = true
	g.bits[2*(x+y*g.width)+1] = v
}

func (g *boolBitGrid) Clear() {
	for i, _ := range g.bits {
		g.bits[i] = false
	}
}
