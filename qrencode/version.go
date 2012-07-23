package qrencode

import (
	"errors"
)

type versionNumber int

type ecbs struct {
	codewordsPerBlock int
	blocks            []ecb
}

type ecb struct {
	count, dataCodewords int
}

func (blocks ecbs) numBlocks() int {
	total := 0
	for _, block := range blocks.blocks {
		total += block.count
	}
	return total
}

func (blocks ecbs) totalECCodewords() int {
	return blocks.numBlocks() * blocks.codewordsPerBlock
}

var (
	ecBlocks = [41][4]ecbs{
		[4]ecbs{},
		[4]ecbs{ // 1
			ecbs{10, []ecb{ecb{1, 16}}}, // M
			ecbs{7, []ecb{ecb{1, 19}}},  // L
			ecbs{17, []ecb{ecb{1, 9}}},  // H
			ecbs{13, []ecb{ecb{1, 13}}}, // Q
		},
		[4]ecbs{ // 2
			ecbs{16, []ecb{ecb{1, 28}}}, // M
			ecbs{10, []ecb{ecb{1, 34}}}, // L
			ecbs{28, []ecb{ecb{1, 16}}}, // H
			ecbs{22, []ecb{ecb{1, 22}}}, // Q
		},
		[4]ecbs{ // 3
			ecbs{26, []ecb{ecb{1, 44}}}, // M
			ecbs{15, []ecb{ecb{1, 55}}}, // L
			ecbs{22, []ecb{ecb{2, 13}}}, // H
			ecbs{18, []ecb{ecb{2, 17}}}, // Q
		},
		[4]ecbs{ // 4
			ecbs{18, []ecb{ecb{2, 32}}}, // M
			ecbs{20, []ecb{ecb{1, 80}}}, // L
			ecbs{16, []ecb{ecb{4, 9}}},  // H
			ecbs{26, []ecb{ecb{2, 24}}}, // Q
		},
		[4]ecbs{ // 5
			ecbs{24, []ecb{ecb{2, 43}}},             // M
			ecbs{26, []ecb{ecb{1, 108}}},            // L
			ecbs{22, []ecb{ecb{2, 11}, ecb{2, 12}}}, // H
			ecbs{18, []ecb{ecb{2, 15}, ecb{2, 16}}}, // Q
		},
		[4]ecbs{ // 6
			ecbs{16, []ecb{ecb{4, 27}}}, // M
			ecbs{18, []ecb{ecb{2, 68}}}, // L
			ecbs{28, []ecb{ecb{4, 15}}}, // H
			ecbs{24, []ecb{ecb{4, 19}}}, // Q
		},
		[4]ecbs{ // 7
			ecbs{18, []ecb{ecb{4, 31}}},             // M
			ecbs{20, []ecb{ecb{2, 78}}},             // L
			ecbs{26, []ecb{ecb{4, 13}, ecb{1, 14}}}, // H
			ecbs{18, []ecb{ecb{2, 14}, ecb{4, 15}}}, // Q
		},
		[4]ecbs{ // 8
			ecbs{22, []ecb{ecb{2, 38}, ecb{2, 39}}}, // M
			ecbs{24, []ecb{ecb{2, 97}}},             // L
			ecbs{26, []ecb{ecb{4, 14}, ecb{2, 15}}}, // H
			ecbs{22, []ecb{ecb{4, 18}, ecb{2, 19}}}, // Q
		},
		[4]ecbs{ // 9
			ecbs{22, []ecb{ecb{3, 36}, ecb{2, 37}}}, // M
			ecbs{30, []ecb{ecb{2, 116}}},            // L
			ecbs{24, []ecb{ecb{4, 12}, ecb{4, 13}}}, // H
			ecbs{20, []ecb{ecb{4, 16}, ecb{4, 17}}}, // Q
		},
		[4]ecbs{ // 10
			ecbs{26, []ecb{ecb{4, 43}, ecb{1, 44}}}, // M
			ecbs{18, []ecb{ecb{2, 68}, ecb{2, 69}}}, // L
			ecbs{28, []ecb{ecb{6, 15}, ecb{2, 16}}}, // H
			ecbs{24, []ecb{ecb{6, 19}, ecb{2, 20}}}, // Q
		},
		[4]ecbs{ // 11
			ecbs{30, []ecb{ecb{1, 50}, ecb{4, 51}}}, // M
			ecbs{20, []ecb{ecb{4, 81}}},             // L
			ecbs{24, []ecb{ecb{3, 12}, ecb{8, 13}}}, // H
			ecbs{28, []ecb{ecb{4, 22}, ecb{4, 23}}}, // Q
		},
		[4]ecbs{ // 12
			ecbs{22, []ecb{ecb{6, 36}, ecb{2, 37}}}, // M
			ecbs{24, []ecb{ecb{2, 92}, ecb{2, 93}}}, // L
			ecbs{28, []ecb{ecb{7, 14}, ecb{4, 15}}}, // H
			ecbs{26, []ecb{ecb{4, 20}, ecb{6, 21}}}, // Q
		},
		[4]ecbs{ // 13
			ecbs{22, []ecb{ecb{8, 37}, ecb{1, 38}}},  // M
			ecbs{26, []ecb{ecb{4, 107}}},             // L
			ecbs{22, []ecb{ecb{12, 11}, ecb{4, 12}}}, // H
			ecbs{24, []ecb{ecb{8, 20}, ecb{4, 21}}},  // Q
		},
		[4]ecbs{ // 14
			ecbs{24, []ecb{ecb{4, 40}, ecb{5, 41}}},   // M
			ecbs{30, []ecb{ecb{3, 115}, ecb{1, 116}}}, // L
			ecbs{24, []ecb{ecb{11, 12}, ecb{5, 13}}},  // H
			ecbs{20, []ecb{ecb{11, 16}, ecb{5, 17}}},  // Q
		},
		[4]ecbs{ // 15
			ecbs{24, []ecb{ecb{5, 41}, ecb{5, 42}}},  // M
			ecbs{22, []ecb{ecb{5, 87}, ecb{1, 88}}},  // L
			ecbs{24, []ecb{ecb{11, 12}, ecb{7, 13}}}, // H
			ecbs{30, []ecb{ecb{5, 24}, ecb{7, 25}}},  // Q
		},
		[4]ecbs{ // 16
			ecbs{28, []ecb{ecb{7, 45}, ecb{3, 46}}},  // M
			ecbs{24, []ecb{ecb{5, 98}, ecb{1, 99}}},  // L
			ecbs{30, []ecb{ecb{3, 15}, ecb{13, 16}}}, // H
			ecbs{24, []ecb{ecb{15, 19}, ecb{2, 20}}}, // Q
		},
		[4]ecbs{ // 17
			ecbs{28, []ecb{ecb{10, 46}, ecb{1, 47}}},  // M
			ecbs{28, []ecb{ecb{1, 107}, ecb{5, 108}}}, // L
			ecbs{28, []ecb{ecb{2, 14}, ecb{17, 15}}},  // H
			ecbs{28, []ecb{ecb{1, 22}, ecb{15, 23}}},  // Q
		},
		[4]ecbs{ // 18
			ecbs{26, []ecb{ecb{9, 43}, ecb{4, 44}}},   // M
			ecbs{30, []ecb{ecb{5, 120}, ecb{1, 121}}}, // L
			ecbs{28, []ecb{ecb{2, 14}, ecb{19, 15}}},  // H
			ecbs{28, []ecb{ecb{17, 22}, ecb{1, 23}}},  // Q
		},
		[4]ecbs{ // 19
			ecbs{26, []ecb{ecb{3, 44}, ecb{11, 44}}},  // M
			ecbs{28, []ecb{ecb{3, 113}, ecb{4, 114}}}, // L
			ecbs{26, []ecb{ecb{9, 13}, ecb{16, 14}}},  // H
			ecbs{26, []ecb{ecb{17, 21}, ecb{4, 22}}},  // Q
		},
		[4]ecbs{ // 20
			ecbs{26, []ecb{ecb{3, 41}, ecb{13, 42}}},  // M
			ecbs{28, []ecb{ecb{3, 107}, ecb{5, 108}}}, // L
			ecbs{28, []ecb{ecb{15, 15}, ecb{10, 16}}}, // H
			ecbs{30, []ecb{ecb{15, 24}, ecb{5, 25}}},  // Q
		},
		[4]ecbs{ // 21
			ecbs{26, []ecb{ecb{17, 42}}},              // M
			ecbs{28, []ecb{ecb{4, 116}, ecb{4, 117}}}, // L
			ecbs{30, []ecb{ecb{19, 16}, ecb{6, 17}}},  // H
			ecbs{28, []ecb{ecb{17, 22}, ecb{6, 23}}},  // Q
		},
		[4]ecbs{ // 22
			ecbs{28, []ecb{ecb{17, 46}}},              // M
			ecbs{28, []ecb{ecb{2, 111}, ecb{7, 112}}}, // L
			ecbs{24, []ecb{ecb{34, 13}}},              // H
			ecbs{30, []ecb{ecb{7, 24}, ecb{16, 25}}},  // Q
		},
		[4]ecbs{ // 23
			ecbs{28, []ecb{ecb{4, 47}, ecb{14, 48}}},  // M
			ecbs{30, []ecb{ecb{4, 121}, ecb{5, 122}}}, // L
			ecbs{30, []ecb{ecb{16, 15}, ecb{14, 16}}}, // H
			ecbs{30, []ecb{ecb{11, 24}, ecb{14, 25}}}, // Q
		},
		[4]ecbs{ // 24
			ecbs{28, []ecb{ecb{6, 45}, ecb{14, 46}}},  // M
			ecbs{30, []ecb{ecb{6, 117}, ecb{4, 118}}}, // L
			ecbs{30, []ecb{ecb{30, 16}, ecb{2, 17}}},  // H
			ecbs{30, []ecb{ecb{11, 24}, ecb{16, 25}}}, // Q
		},
		[4]ecbs{ // 25
			ecbs{28, []ecb{ecb{8, 47}, ecb{13, 48}}},  // M
			ecbs{26, []ecb{ecb{8, 106}, ecb{4, 107}}}, // L
			ecbs{30, []ecb{ecb{22, 15}, ecb{13, 16}}}, // H
			ecbs{30, []ecb{ecb{7, 24}, ecb{22, 25}}},  // Q
		},
		[4]ecbs{ // 26
			ecbs{28, []ecb{ecb{19, 46}, ecb{4, 47}}},   // M
			ecbs{28, []ecb{ecb{10, 114}, ecb{2, 115}}}, // L
			ecbs{30, []ecb{ecb{33, 16}, ecb{4, 17}}},   // H
			ecbs{28, []ecb{ecb{28, 22}, ecb{6, 23}}},   // Q
		},
		[4]ecbs{ // 27
			ecbs{28, []ecb{ecb{22, 45}, ecb{3, 46}}},  // M
			ecbs{30, []ecb{ecb{8, 122}, ecb{4, 123}}}, // L
			ecbs{30, []ecb{ecb{12, 15}, ecb{28, 16}}}, // H
			ecbs{30, []ecb{ecb{8, 23}, ecb{26, 24}}},  // Q
		},
		[4]ecbs{ // 28
			ecbs{28, []ecb{ecb{3, 45}, ecb{23, 46}}},   // M
			ecbs{30, []ecb{ecb{3, 117}, ecb{10, 118}}}, // L
			ecbs{30, []ecb{ecb{11, 15}, ecb{31, 16}}},  // H
			ecbs{30, []ecb{ecb{4, 24}, ecb{31, 25}}},   // Q
		},
		[4]ecbs{ // 29
			ecbs{28, []ecb{ecb{21, 45}, ecb{7, 46}}},  // M
			ecbs{30, []ecb{ecb{7, 116}, ecb{7, 117}}}, // L
			ecbs{30, []ecb{ecb{19, 15}, ecb{26, 16}}}, // H
			ecbs{30, []ecb{ecb{1, 32}, ecb{37, 24}}},  // Q
		},
		[4]ecbs{ // 30
			ecbs{28, []ecb{ecb{19, 47}, ecb{10, 48}}},  // M
			ecbs{30, []ecb{ecb{5, 115}, ecb{10, 116}}}, // L
			ecbs{30, []ecb{ecb{23, 15}, ecb{25, 16}}},  // H
			ecbs{30, []ecb{ecb{15, 24}, ecb{25, 25}}},  // Q
		},
		[4]ecbs{ // 31
			ecbs{28, []ecb{ecb{2, 46}, ecb{29, 47}}},   // M
			ecbs{30, []ecb{ecb{13, 115}, ecb{3, 116}}}, // L
			ecbs{30, []ecb{ecb{23, 15}, ecb{28, 16}}},  // H
			ecbs{30, []ecb{ecb{42, 24}, ecb{1, 25}}},   // Q
		},
		[4]ecbs{ // 32
			ecbs{28, []ecb{ecb{10, 46}, ecb{23, 47}}}, // M
			ecbs{30, []ecb{ecb{17, 115}}},             // L
			ecbs{30, []ecb{ecb{19, 15}, ecb{35, 16}}}, // H
			ecbs{30, []ecb{ecb{10, 24}, ecb{35, 25}}}, // Q
		},
		[4]ecbs{ // 33
			ecbs{28, []ecb{ecb{14, 46}, ecb{21, 47}}},  // M
			ecbs{30, []ecb{ecb{17, 115}, ecb{1, 116}}}, // L
			ecbs{30, []ecb{ecb{11, 15}, ecb{46, 16}}},  // H
			ecbs{30, []ecb{ecb{29, 24}, ecb{19, 25}}},  // Q
		},
		[4]ecbs{ // 34
			ecbs{28, []ecb{ecb{14, 16}, ecb{23, 47}}},  // M
			ecbs{30, []ecb{ecb{13, 115}, ecb{6, 116}}}, // L
			ecbs{30, []ecb{ecb{59, 16}, ecb{1, 17}}},   // H
			ecbs{30, []ecb{ecb{44, 24}, ecb{7, 25}}},   // Q
		},
		[4]ecbs{ // 35
			ecbs{28, []ecb{ecb{12, 47}, ecb{26, 48}}},  // M
			ecbs{30, []ecb{ecb{12, 121}, ecb{7, 122}}}, // L
			ecbs{30, []ecb{ecb{22, 15}, ecb{41, 16}}},  // H
			ecbs{30, []ecb{ecb{39, 24}, ecb{14, 25}}},  // Q
		},
		[4]ecbs{ // 36
			ecbs{28, []ecb{ecb{6, 47}, ecb{34, 48}}},   // M
			ecbs{30, []ecb{ecb{6, 121}, ecb{14, 122}}}, // L
			ecbs{30, []ecb{ecb{2, 15}, ecb{64, 16}}},   // H
			ecbs{30, []ecb{ecb{46, 24}, ecb{10, 25}}},  // Q
		},
		[4]ecbs{ // 37
			ecbs{28, []ecb{ecb{29, 46}, ecb{14, 47}}},  // M
			ecbs{30, []ecb{ecb{17, 122}, ecb{4, 123}}}, // L
			ecbs{30, []ecb{ecb{24, 15}, ecb{46, 16}}},  // H
			ecbs{30, []ecb{ecb{49, 24}, ecb{10, 25}}},  // Q
		},
		[4]ecbs{ // 38
			ecbs{28, []ecb{ecb{13, 46}, ecb{32, 47}}},  // M
			ecbs{30, []ecb{ecb{4, 122}, ecb{18, 123}}}, // L
			ecbs{30, []ecb{ecb{42, 15}, ecb{32, 16}}},  // H
			ecbs{30, []ecb{ecb{48, 24}, ecb{14, 25}}},  // Q
		},
		[4]ecbs{ // 39
			ecbs{28, []ecb{ecb{40, 47}, ecb{7, 48}}},   // M
			ecbs{30, []ecb{ecb{20, 117}, ecb{4, 118}}}, // L
			ecbs{30, []ecb{ecb{10, 15}, ecb{67, 16}}},  // H
			ecbs{30, []ecb{ecb{43, 24}, ecb{22, 25}}},  // Q
		},
		[4]ecbs{ // 40
			ecbs{28, []ecb{ecb{18, 47}, ecb{31, 48}}},  // M
			ecbs{30, []ecb{ecb{19, 118}, ecb{6, 119}}}, // L
			ecbs{30, []ecb{ecb{20, 15}, ecb{61, 16}}},  // H
			ecbs{30, []ecb{ecb{34, 24}, ecb{34, 25}}},  // Q
		},
	}
)

func (v versionNumber) totalCodewords() int {
	total := 0
	ecCodewords := ecBlocks[v][1].codewordsPerBlock
	for _, block := range ecBlocks[v][1].blocks {
		total += block.count * (block.dataCodewords + ecCodewords)
	}
	return total
}

func (v versionNumber) dimension() int {
	return 17 + 4*int(v)
}

func chooseVersion(bitCount int, ecLevel ECLevel) (versionNumber, error) {
	for v := versionNumber(1); v <= 40; v++ {
		if (bitCount+7)/8 <= v.totalCodewords()-ecBlocks[v][ecLevel].totalECCodewords() {
			return v, nil
		}
	}
	return versionNumber(0), errors.New("Content too large")
}
