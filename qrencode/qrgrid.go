package qrencode

var (
	positionDetectionPatternBack = [][]byte{
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
	}

	positionDetectionPattern = [][]byte{
		{1, 1, 1, 1, 1, 1, 1},
		{1, 0, 0, 0, 0, 0, 1},
		{1, 0, 1, 1, 1, 0, 1},
		{1, 0, 1, 1, 1, 0, 1},
		{1, 0, 1, 1, 1, 0, 1},
		{1, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1},
	}

	positionAdjustmentPattern = [][]byte{
		{1, 1, 1, 1, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 1, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 1, 1, 1, 1},
	}

	positionAdjustmentPatternCoordinateTable = [][]int{
		{},
		{},                             // Version 1
		{6, 18},                        // Version 2
		{6, 22},                        // Version 3
		{6, 26},                        // Version 4
		{6, 30},                        // Version 5
		{6, 34},                        // Version 6
		{6, 22, 38},                    // Version 7
		{6, 24, 42},                    // Version 8
		{6, 26, 46},                    // Version 9
		{6, 28, 50},                    // Version 10
		{6, 30, 54},                    // Version 11
		{6, 32, 58},                    // Version 12
		{6, 34, 62},                    // Version 13
		{6, 26, 46, 66},                // Version 14
		{6, 26, 48, 70},                // Version 15
		{6, 26, 50, 74},                // Version 16
		{6, 30, 54, 78},                // Version 17
		{6, 30, 56, 82},                // Version 18
		{6, 30, 58, 86},                // Version 19
		{6, 34, 62, 90},                // Version 20
		{6, 28, 50, 72, 94},            // Version 21
		{6, 26, 50, 74, 98},            // Version 22
		{6, 30, 54, 78, 102},           // Version 23
		{6, 28, 54, 80, 106},           // Version 24
		{6, 32, 58, 84, 110},           // Version 25
		{6, 30, 58, 86, 114},           // Version 26
		{6, 34, 62, 90, 118},           // Version 27
		{6, 26, 50, 74, 98, 122},       // Version 28
		{6, 30, 54, 78, 102, 126},      // Version 29
		{6, 26, 52, 78, 104, 130},      // Version 30
		{6, 30, 56, 82, 108, 134},      // Version 31
		{6, 34, 60, 86, 112, 138},      // Version 32
		{6, 30, 58, 86, 114, 142},      // Version 33
		{6, 34, 62, 90, 118, 146},      // Version 34
		{6, 30, 54, 78, 102, 126, 150}, // Version 35
		{6, 24, 50, 76, 102, 128, 154}, // Version 36
		{6, 28, 54, 80, 106, 132, 158}, // Version 37
		{6, 32, 58, 84, 110, 136, 162}, // Version 38
		{6, 26, 54, 82, 110, 138, 166}, // Version 39
		{6, 30, 58, 86, 114, 142, 170}, // Version 40
	}

	typeInfoCoordinates = [][]int{
		{8, 0},
		{8, 1},
		{8, 2},
		{8, 3},
		{8, 4},
		{8, 5},
		{8, 7},
		{8, 8},
		{7, 8},
		{5, 8},
		{4, 8},
		{3, 8},
		{2, 8},
		{1, 8},
		{0, 8},
	}
)

const (
	versionInfoPoly     = 0x1f25
	typeInfoPoly        = 0x537
	typeInfoMaskPattern = 0x5412
)

func bestMaskPattern(bits *BitVector, version versionNumber, ecLevel ECLevel, grid *BitGrid) int {
	bestMaskPattern := -1
	bestPenalty := 0
	for maskPattern := 0; maskPattern < 8; maskPattern++ {
		buildGrid(bits, version, ecLevel, maskPattern, grid)
		penalty := maskPenalty(grid)
		if bestMaskPattern < 0 || penalty < bestPenalty {
			bestMaskPattern = maskPattern
			bestPenalty = penalty
		}
	}
	if bestMaskPattern < 0 {
		panic("bestMaskPattern < 0")
	}
	return bestMaskPattern
}

func buildGrid(bits *BitVector, version versionNumber, ecLevel ECLevel, maskPattern int, grid *BitGrid) {
	grid.Clear()
	embedBasicPatterns(version, grid)
	embedTypeInfo(ecLevel, maskPattern, grid)
	maybeEmbedVersionInfo(version, grid)
	embedDataBits(bits, maskPattern, grid)
}

func embedBasicPatterns(version versionNumber, grid *BitGrid) {
	embedPositionDetectionPatternsAndSeparators(grid)
	embedDarkDotAtLeftBottomCorner(grid)
	maybeEmbedPositionAdjustmentPatterns(version, grid)
	embedTimingPatterns(grid)
}

func embedPositionDetectionPatternsAndSeparators(grid *BitGrid) {
	embedPattern(0, 0, positionDetectionPatternBack, grid)
	embedPattern(grid.Width()-len(positionDetectionPatternBack[0]), 0, positionDetectionPatternBack, grid)
	embedPattern(0, grid.Height()-len(positionDetectionPatternBack), positionDetectionPatternBack, grid)
	embedPattern(0, 0, positionDetectionPattern, grid)
	embedPattern(grid.Width()-len(positionDetectionPattern[0]), 0, positionDetectionPattern, grid)
	embedPattern(0, grid.Height()-len(positionDetectionPattern), positionDetectionPattern, grid)
}

func embedDarkDotAtLeftBottomCorner(grid *BitGrid) {
	grid.Set(8, grid.Height()-8, true)
}

func maybeEmbedPositionAdjustmentPatterns(version versionNumber, grid *BitGrid) {
	h := len(positionAdjustmentPattern)
	w := len(positionAdjustmentPattern[h/2])
	for _, y := range positionAdjustmentPatternCoordinateTable[version] {
		for _, x := range positionAdjustmentPatternCoordinateTable[version] {
			if grid.Empty(x, y) {
				embedPattern(x-w/2, y-h/2, positionAdjustmentPattern, grid)
			}
		}
	}
}

func embedPattern(x0, y0 int, pattern [][]byte, grid *BitGrid) {
	for y, row := range pattern {
		for x, v := range row {
			grid.Set(x0+x, y0+y, v != 0)
		}
	}
}

func embedTimingPatterns(grid *BitGrid) {
	for i := 0; i < grid.Width(); i++ {
		if grid.Empty(i, 6) {
			grid.Set(i, 6, i&1 == 0)
		}
	}
	for i := 0; i < grid.Height(); i++ {
		if grid.Empty(6, i) {
			grid.Set(6, i, i&1 == 0)
		}
	}
}

func embedTypeInfo(ecLevel ECLevel, maskPattern int, grid *BitGrid) {
	typeInfo := int(ecLevel)<<3 | maskPattern
	bchCode := calculateBCHCode(typeInfo, typeInfoPoly)
	typeInfo = (typeInfo<<10 | (bchCode & 0x3ff)) ^ typeInfoMaskPattern
	for i := 0; i < 15; i++ {
		bit := (typeInfo>>uint(i))&1 == 1
		grid.Set(typeInfoCoordinates[i][0], typeInfoCoordinates[i][1], bit)
		if i < 8 {
			grid.Set(grid.Width()-i-1, 8, bit)
		} else {
			grid.Set(8, grid.Height()+i-15, bit)
		}
	}
}

func maybeEmbedVersionInfo(version versionNumber, grid *BitGrid) {
	if version < 7 {
		return
	}
	versionInfo := int(version)<<12 | calculateBCHCode(int(version), versionInfoPoly)
	for i := 0; i < 6; i++ {
		for j := 0; j < 3; j++ {
			bit := versionInfo&1 == 1
			versionInfo >>= 1
			grid.Set(i, grid.Height()-11+j, bit)
			grid.Set(grid.Width()-11+j, i, bit)
		}
	}
}

func embedDataBits(bits *BitVector, maskPattern int, grid *BitGrid) {
	bitIndex := 0
	for direction, x, y := -1, grid.Width()-1, grid.Height()-1; x > 0; x, y, direction = x-2, y-direction, -direction {
		if x == 6 {
			x--
		}
		for ; y >= 0 && y < grid.Height(); y += direction {
			for i := 0; i < 2; i++ {
				xx := x - i
				if !grid.Empty(xx, y) {
					continue
				}
				bit := false
				if bitIndex < bits.Length() {
					bit = bits.Get(bitIndex)
					bitIndex++
				}
				if mask(maskPattern, xx, y) {
					bit = !bit
				}
				grid.Set(xx, y, bit)
			}
		}
	}
	if bitIndex != bits.Length() {
		panic("bitIndex != bits.Length()")
	}
}

func calculateBCHCode(value, poly int) int {
	msbSetInPoly := findMSBSet(poly)
	value <<= uint(msbSetInPoly - 1)
	for findMSBSet(value) >= msbSetInPoly {
		value ^= poly << uint(findMSBSet(value)-msbSetInPoly)
	}
	return value
}

func findMSBSet(value int) int {
	numDigits := 0
	for v := uint(value); v != 0; v >>= 1 {
		numDigits++
	}
	return numDigits
}

func mask(maskPattern int, x, y int) bool {
	switch maskPattern {
	case 0:
		return (x+y)&1 == 0
	case 1:
		return y&1 == 0
	case 2:
		return x%3 == 0
	case 3:
		return (x+y)%3 == 0
	case 4:
		return (x/3+y>>1)&1 == 0
	case 5:
		return (x*y)&1+(x*y)%3 == 0
	case 6:
		return ((x*y)&1+(x*y)%3)&1 == 0
	case 7:
		return ((x*y)%3+(x+y)&1)&1 == 0
	}
	return false
}

func maskPenalty(grid *BitGrid) int {
	return maskPenaltyRule1(grid) + maskPenaltyRule2(grid) + maskPenaltyRule3(grid) + maskPenaltyRule4(grid)
}

func maskPenaltyRule1(grid *BitGrid) int {
	penalty := 0
	for x := 0; x < grid.Width(); x++ {
		for y, count := 1, 1; y < grid.Height(); y++ {
			if grid.Get(x, y) == grid.Get(x, y-1) {
				count++
				if count == 5 {
					penalty += 3
				} else if count > 5 {
					penalty++
				}
			} else {
				count = 1
			}
		}
	}
	for y := 0; y < grid.Height(); y++ {
		for x, count := 1, 1; x < grid.Width(); x++ {
			if grid.Get(x, y) == grid.Get(x-1, y) {
				count++
				if count == 5 {
					penalty += 3
				} else if count > 5 {
					penalty++
				}
			} else {
				count = 1
			}
		}
	}
	return penalty
}

func maskPenaltyRule2(grid *BitGrid) int {
	penalty := 0
	for y := 1; y < grid.Height(); y++ {
		for x := 1; x < grid.Width(); x++ {
			v := grid.Get(x, y)
			if v == grid.Get(x-1, y) && v == grid.Get(x-1, y-1) && v == grid.Get(x, y-1) {
				penalty += 3
			}
		}
	}
	return penalty
}

func maskPenaltyRule3(grid *BitGrid) int {
	penalty := 0
	for y := 0; y < grid.Height(); y++ {
		for x := 0; x < grid.Width(); x++ {
			if x+6 < grid.Width() &&
				grid.Get(x, y) &&
				!grid.Get(x+1, y) &&
				grid.Get(x+2, y) &&
				grid.Get(x+3, y) &&
				grid.Get(x+4, y) &&
				!grid.Get(x+5, y) &&
				grid.Get(x+6, y) {
				if x+10 < grid.Width() &&
					!grid.Get(x+7, y) &&
					!grid.Get(x+8, y) &&
					!grid.Get(x+9, y) &&
					!grid.Get(x+10, y) {
					penalty += 40
				}
				if x-4 >= 0 &&
					!grid.Get(x-4, y) &&
					!grid.Get(x-3, y) &&
					!grid.Get(x-2, y) &&
					!grid.Get(x-1, y) {
					penalty += 40
				}
			}
			if y+6 < grid.Height() &&
				grid.Get(x, y) &&
				!grid.Get(x, y+1) &&
				grid.Get(x, y+2) &&
				grid.Get(x, y+3) &&
				grid.Get(x, y+4) &&
				!grid.Get(x, y+5) &&
				grid.Get(x, y+6) {
				if y+10 < grid.Height() &&
					!grid.Get(x, y+7) &&
					!grid.Get(x, y+8) &&
					!grid.Get(x, y+9) &&
					!grid.Get(x, y+10) {
					penalty += 40
				}
				if y-4 >= 0 &&
					!grid.Get(x, y-4) &&
					!grid.Get(x, y-3) &&
					!grid.Get(x, y-2) &&
					!grid.Get(x, y-1) {
					penalty += 40
				}
			}
		}
	}
	return penalty
}

func maskPenaltyRule4(grid *BitGrid) int {
	total := grid.Width() * grid.Height()
	dark := -total / 2
	for y := 0; y < grid.Height(); y++ {
		for x := 0; x < grid.Width(); x++ {
			if grid.Get(x, y) {
				dark++
			}
		}
	}
	if dark < 0 {
		dark = 1 - dark
	}
	return (dark * 200) / total
}
