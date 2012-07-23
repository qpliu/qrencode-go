package qrencode

type ECLevel int

const (
	ECLevelM = ECLevel(0)
	ECLevelL = ECLevel(1)
	ECLevelH = ECLevel(2)
	ECLevelQ = ECLevel(3)
)

func supportedECLevel(ecLevel ECLevel) bool {
	switch ecLevel {
	case ECLevelM, ECLevelL, ECLevelH, ECLevelQ:
		return true
	}
	return false
}

type blockPair struct {
	dataBytes, ecBytes []int
}

func interleaveWithECBytes(bits *BitVector, version versionNumber, ecLevel ECLevel) *BitVector {
	numTotalBytes := version.totalCodewords()
	numDataBytes := version.totalCodewords() - ecBlocks[version][ecLevel].totalECCodewords()
	numRSBlocks := ecBlocks[version][ecLevel].numBlocks()
	if bits.Length() != numDataBytes*8 {
		panic("bits.Length() != numDataBytes*8")
	}

	dataBytesOffset := 0
	maxNumDataBytes := 0
	maxNumEcBytes := 0

	blocks := make([]blockPair, numRSBlocks)

	for i := 0; i < numRSBlocks; i++ {
		numDataBytes, numEcBytes := getBlockSizes(numTotalBytes, numDataBytes, numRSBlocks, i)
		blocks[i] = blockPair{make([]int, numDataBytes), make([]int, numEcBytes)}
		for j := 0; j < numDataBytes; j++ {
			blocks[i].dataBytes[j] = 0
			if bits.Get(8 * (dataBytesOffset + j)) {
				blocks[i].dataBytes[j] |= 128
			}
			if bits.Get(8*(dataBytesOffset+j) + 1) {
				blocks[i].dataBytes[j] |= 64
			}
			if bits.Get(8*(dataBytesOffset+j) + 2) {
				blocks[i].dataBytes[j] |= 32
			}
			if bits.Get(8*(dataBytesOffset+j) + 3) {
				blocks[i].dataBytes[j] |= 16
			}
			if bits.Get(8*(dataBytesOffset+j) + 4) {
				blocks[i].dataBytes[j] |= 8
			}
			if bits.Get(8*(dataBytesOffset+j) + 5) {
				blocks[i].dataBytes[j] |= 4
			}
			if bits.Get(8*(dataBytesOffset+j) + 6) {
				blocks[i].dataBytes[j] |= 2
			}
			if bits.Get(8*(dataBytesOffset+j) + 7) {
				blocks[i].dataBytes[j] |= 1
			}
		}
		generateECBytes(&blocks[i])

		if numDataBytes > maxNumDataBytes {
			maxNumDataBytes = numDataBytes
		}
		if numEcBytes > maxNumEcBytes {
			maxNumEcBytes = numEcBytes
		}
		dataBytesOffset += numDataBytes
	}

	if numDataBytes != dataBytesOffset {
		panic("numDataBytes != dataBytesOffset")
	}

	result := &BitVector{}
	for i := 0; i < maxNumDataBytes; i++ {
		for _, block := range blocks {
			if i < len(block.dataBytes) {
				result.Append(block.dataBytes[i], 8)
			}
		}
	}
	for i := 0; i < maxNumEcBytes; i++ {
		for _, block := range blocks {
			if i < len(block.ecBytes) {
				result.Append(block.ecBytes[i], 8)
			}
		}
	}

	if result.Length() != numTotalBytes*8 {
		panic("result.Length() != numTotalBytes*8")
	}
	return result
}

func getBlockSizes(numTotalBytes, numDataBytes, numRSBlocks, blockID int) (int, int) {
	if blockID >= numRSBlocks {
		panic("blockID >= numRSBlocks")
	}
	numRsBlocksInGroup2 := numTotalBytes % numRSBlocks
	numRsBlocksInGroup1 := numRSBlocks - numRsBlocksInGroup2
	numTotalBytesInGroup1 := numTotalBytes / numRSBlocks
	numTotalBytesInGroup2 := numTotalBytesInGroup1 + 1
	numDataBytesInGroup1 := numDataBytes / numRSBlocks
	numDataBytesInGroup2 := numDataBytesInGroup1 + 1
	numEcBytesInGroup1 := numTotalBytesInGroup1 - numDataBytesInGroup1
	numEcBytesInGroup2 := numTotalBytesInGroup2 - numDataBytesInGroup2
	if numEcBytesInGroup1 != numEcBytesInGroup2 {
		panic("numEcBytesInGroup1 != numEcBytesInGroup2")
	}
	if numRSBlocks != numRsBlocksInGroup1+numRsBlocksInGroup2 {
		panic("numRSBlocks != numRsBlocksInGroup1 + numRsBlocksInGroup2")
	}
	if numTotalBytes != (numDataBytesInGroup1+numEcBytesInGroup1)*numRsBlocksInGroup1+(numDataBytesInGroup2+numEcBytesInGroup2)*numRsBlocksInGroup2 {
		panic("numTotalBytes != (numDataBytesInGroup1 + numEcBytesInGroup1)*numRsBlocksInGroup1 + (numDataBytesInGroup2 + numEcBytesInGroup2)*numRsBlocksInGroup2")
	}
	if blockID < numRsBlocksInGroup1 {
		return numDataBytesInGroup1, numEcBytesInGroup1
	}
	return numDataBytesInGroup2, numEcBytesInGroup2
}

func generateECBytes(block *blockPair) {
	generator := buildGenerator(len(block.ecBytes))
	info := newGFPoly(block.dataBytes)
	info = info.MultiplyByMonomial(len(block.ecBytes), 1)
	_, remainder := info.Divide(generator)
	numZeroCoefficients := len(block.ecBytes) - len(remainder.coefficients)
	for i := 0; i < numZeroCoefficients; i++ {
		block.ecBytes[i] = 0
	}
	copy(block.ecBytes[numZeroCoefficients:], remainder.coefficients)
}

var (
	fieldExpTable = [256]int{
		1, 2, 4, 8, 16, 32, 64, 128, 29, 58, 116, 232, 205, 135, 19, 38,
		76, 152, 45, 90, 180, 117, 234, 201, 143, 3, 6, 12, 24, 48, 96, 192,
		157, 39, 78, 156, 37, 74, 148, 53, 106, 212, 181, 119, 238, 193, 159, 35,
		70, 140, 5, 10, 20, 40, 80, 160, 93, 186, 105, 210, 185, 111, 222, 161,
		95, 190, 97, 194, 153, 47, 94, 188, 101, 202, 137, 15, 30, 60, 120, 240,
		253, 231, 211, 187, 107, 214, 177, 127, 254, 225, 223, 163, 91, 182, 113, 226,
		217, 175, 67, 134, 17, 34, 68, 136, 13, 26, 52, 104, 208, 189, 103, 206,
		129, 31, 62, 124, 248, 237, 199, 147, 59, 118, 236, 197, 151, 51, 102, 204,
		133, 23, 46, 92, 184, 109, 218, 169, 79, 158, 33, 66, 132, 21, 42, 84,
		168, 77, 154, 41, 82, 164, 85, 170, 73, 146, 57, 114, 228, 213, 183, 115,
		230, 209, 191, 99, 198, 145, 63, 126, 252, 229, 215, 179, 123, 246, 241, 255,
		227, 219, 171, 75, 150, 49, 98, 196, 149, 55, 110, 220, 165, 87, 174, 65,
		130, 25, 50, 100, 200, 141, 7, 14, 28, 56, 112, 224, 221, 167, 83, 166,
		81, 162, 89, 178, 121, 242, 249, 239, 195, 155, 43, 86, 172, 69, 138, 9,
		18, 36, 72, 144, 61, 122, 244, 245, 247, 243, 251, 235, 203, 139, 11, 22,
		44, 88, 176, 125, 250, 233, 207, 131, 27, 54, 108, 216, 173, 71, 142, 1,
	}
	fieldLogTable = [256]int{
		0, 0, 1, 25, 2, 50, 26, 198, 3, 223, 51, 238, 27, 104, 199, 75,
		4, 100, 224, 14, 52, 141, 239, 129, 28, 193, 105, 248, 200, 8, 76, 113,
		5, 138, 101, 47, 225, 36, 15, 33, 53, 147, 142, 218, 240, 18, 130, 69,
		29, 181, 194, 125, 106, 39, 249, 185, 201, 154, 9, 120, 77, 228, 114, 166,
		6, 191, 139, 98, 102, 221, 48, 253, 226, 152, 37, 179, 16, 145, 34, 136,
		54, 208, 148, 206, 143, 150, 219, 189, 241, 210, 19, 92, 131, 56, 70, 64,
		30, 66, 182, 163, 195, 72, 126, 110, 107, 58, 40, 84, 250, 133, 186, 61,
		202, 94, 155, 159, 10, 21, 121, 43, 78, 212, 229, 172, 115, 243, 167, 87,
		7, 112, 192, 247, 140, 128, 99, 13, 103, 74, 222, 237, 49, 197, 254, 24,
		227, 165, 153, 119, 38, 184, 180, 124, 17, 68, 146, 217, 35, 32, 137, 46,
		55, 63, 209, 91, 149, 188, 207, 205, 144, 135, 151, 178, 220, 252, 190, 97,
		242, 86, 211, 171, 20, 42, 93, 158, 132, 60, 57, 83, 71, 109, 65, 162,
		31, 45, 67, 216, 183, 123, 164, 118, 196, 23, 73, 236, 127, 12, 111, 246,
		108, 161, 59, 82, 41, 157, 85, 170, 251, 96, 134, 177, 187, 204, 62, 90,
		203, 89, 95, 176, 156, 169, 160, 81, 11, 245, 22, 235, 122, 117, 44, 215,
		79, 174, 213, 233, 230, 231, 173, 232, 116, 214, 244, 234, 168, 80, 88, 175,
	}
	fieldZero        = gfPoly{[]int{0}}
	cachedGenerators = []gfPoly{gfPoly{[]int{1}}}
)

func fieldAddSub(a, b int) int {
	return a ^ b
}

func fieldInv(a int) int {
	if a == 0 {
		panic("a == 0")
	}
	return fieldExpTable[256-fieldLogTable[a]-1]
}

func fieldMult(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
	return fieldExpTable[(fieldLogTable[a]+fieldLogTable[b])%255]
}

func fieldBuildMonomial(degree, coefficient int) gfPoly {
	if degree < 0 {
		panic("degree < 0")
	}
	if coefficient == 0 {
		return fieldZero
	}
	coeff := make([]int, degree+1)
	coeff[0] = coefficient
	return newGFPoly(coeff)
}

type gfPoly struct {
	coefficients []int
}

func buildGenerator(degree int) gfPoly {
	for degree >= len(cachedGenerators) {
		d := len(cachedGenerators) - 1
		lastGenerator := cachedGenerators[d]
		cachedGenerators = append(cachedGenerators, lastGenerator.Multiply(newGFPoly([]int{1, fieldExpTable[d]})))
	}
	return cachedGenerators[degree]
}

func newGFPoly(coefficients []int) gfPoly {
	if len(coefficients) == 0 {
		panic("len(coefficients) == 0")
	}
	for len(coefficients) > 0 {
		if coefficients[0] != 0 {
			return gfPoly{coefficients}
		}
		coefficients = coefficients[1:]
	}
	return fieldZero
}

func (p gfPoly) Degree() int {
	return len(p.coefficients) - 1
}

func (p gfPoly) IsZero() bool {
	return p.coefficients[0] == 0
}

func (p gfPoly) AddSub(x gfPoly) gfPoly {
	if p.IsZero() {
		return x
	}
	if x.IsZero() {
		return p
	}
	smallerCoeff := p.coefficients
	largerCoeff := x.coefficients
	if len(smallerCoeff) > len(largerCoeff) {
		smallerCoeff, largerCoeff = largerCoeff, smallerCoeff
	}
	coeff := make([]int, len(largerCoeff))
	lenDiff := len(largerCoeff) - len(smallerCoeff)
	copy(coeff, largerCoeff[:lenDiff])
	for i := lenDiff; i < len(coeff); i++ {
		coeff[i] = fieldAddSub(smallerCoeff[i-lenDiff], largerCoeff[i])
	}
	return newGFPoly(coeff)
}

func (p gfPoly) Multiply(x gfPoly) gfPoly {
	if p.IsZero() || x.IsZero() {
		return fieldZero
	}
	coeff := make([]int, len(x.coefficients)+len(p.coefficients)-1)
	for i, a := range p.coefficients {
		for j, b := range x.coefficients {
			coeff[i+j] = fieldAddSub(coeff[i+j], fieldMult(a, b))
		}
	}
	return newGFPoly(coeff)
}

func (p gfPoly) MultiplyByMonomial(degree, coefficient int) gfPoly {
	if degree < 0 {
		panic("degree < 0")
	}
	if coefficient == 0 {
		return fieldZero
	}
	coeff := make([]int, len(p.coefficients)+degree)
	for i, c := range p.coefficients {
		coeff[i] = fieldMult(c, coefficient)
	}
	return newGFPoly(coeff)
}

func (p gfPoly) Divide(x gfPoly) (gfPoly, gfPoly) {
	if x.IsZero() {
		panic("x.IsZero()")
	}
	quotient := fieldZero
	remainder := p

	inverseDenominatorLeadingTerm := fieldInv(x.coefficients[0])

	for remainder.Degree() >= x.Degree() && !remainder.IsZero() {
		degreeDifference := remainder.Degree() - x.Degree()
		scale := fieldMult(remainder.coefficients[0], inverseDenominatorLeadingTerm)
		term := x.MultiplyByMonomial(degreeDifference, scale)
		iterationQuotient := fieldBuildMonomial(degreeDifference, scale)
		quotient = quotient.AddSub(iterationQuotient)
		remainder = remainder.AddSub(term)
	}
	return quotient, remainder
}
