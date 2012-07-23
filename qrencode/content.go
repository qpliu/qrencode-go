package qrencode

import (
	"errors"
)

func stringContentBits(content string, ecLevel ECLevel) (*BitVector, versionNumber, error) {
	if !supportedECLevel(ecLevel) {
		return nil, versionNumber(0), errors.New("Unrecognized ECLevel")
	}
	headerBits := &BitVector{}
	mode := getMode(content)
	if mode == modeByte {
		headerBits.Append(int(modeECI), 4)
		headerBits.Append(26, 8) // UTF-8
	}
	headerBits.Append(int(mode), 4)
	return contentBits([]byte(content), ecLevel, mode, headerBits)
}

func binaryContentBits(content []byte, ecLevel ECLevel) (*BitVector, versionNumber, error) {
	if !supportedECLevel(ecLevel) {
		return nil, versionNumber(0), errors.New("Unrecognized ECLevel")
	}
	headerBits := &BitVector{}
	headerBits.Append(int(modeByte), 4)
	return contentBits(content, ecLevel, modeByte, headerBits)
}

func contentBits(content []byte, ecLevel ECLevel, mode modeIndicator, headerBits *BitVector) (*BitVector, versionNumber, error) {
	dataBits := BitVector{}
	appendContent(content, mode, &dataBits)

	bitsNeeded := headerBits.Length() + dataBits.Length() + mode.characterCountBits(versionNumber(40))
	version, err := chooseVersion(bitsNeeded, ecLevel)
	if err != nil {
		return nil, version, err
	}

	headerAndDataBits := &BitVector{}
	headerAndDataBits.AppendBits(*headerBits)
	headerAndDataBits.Append(len(content), mode.characterCountBits(version))
	headerAndDataBits.AppendBits(dataBits)

	appendTerminator(version.totalCodewords()-ecBlocks[version][ecLevel].totalECCodewords(), headerAndDataBits)
	return headerAndDataBits, version, nil
}

var (
	invalidAlphanumericByte = errors.New("Invalid Alphanumeric Byte")
)

func alphanumericCode(b byte) (byte, error) {
	switch b {
	case 0x30:
		return 0, nil
	case 0x31:
		return 1, nil
	case 0x32:
		return 2, nil
	case 0x33:
		return 3, nil
	case 0x34:
		return 4, nil
	case 0x35:
		return 5, nil
	case 0x36:
		return 6, nil
	case 0x37:
		return 7, nil
	case 0x38:
		return 8, nil
	case 0x39:
		return 9, nil
	case 0x41:
		return 10, nil
	case 0x42:
		return 11, nil
	case 0x43:
		return 12, nil
	case 0x44:
		return 13, nil
	case 0x45:
		return 14, nil
	case 0x46:
		return 15, nil
	case 0x47:
		return 16, nil
	case 0x48:
		return 17, nil
	case 0x49:
		return 18, nil
	case 0x4a:
		return 19, nil
	case 0x4b:
		return 20, nil
	case 0x4c:
		return 21, nil
	case 0x4d:
		return 22, nil
	case 0x4e:
		return 23, nil
	case 0x4f:
		return 24, nil
	case 0x50:
		return 25, nil
	case 0x51:
		return 26, nil
	case 0x52:
		return 27, nil
	case 0x53:
		return 28, nil
	case 0x54:
		return 29, nil
	case 0x55:
		return 30, nil
	case 0x56:
		return 31, nil
	case 0x57:
		return 32, nil
	case 0x58:
		return 33, nil
	case 0x59:
		return 34, nil
	case 0x5a:
		return 35, nil
	case 0x20:
		return 36, nil
	case 0x24:
		return 37, nil
	case 0x25:
		return 38, nil
	case 0x2a:
		return 39, nil
	case 0x2b:
		return 40, nil
	case 0x2d:
		return 41, nil
	case 0x2e:
		return 42, nil
	case 0x2f:
		return 43, nil
	case 0x3a:
		return 44, nil
	}
	return 0, invalidAlphanumericByte
}

func appendContent(content []byte, mode modeIndicator, bits *BitVector) {
	switch mode {
	case modeNumeric:
		for i := 0; i+2 < len(content); i += 3 {
			n1, err := alphanumericCode(content[i])
			if n1 > 9 || err != nil {
				panic("Invalid numeric mode content")
			}
			n2, err := alphanumericCode(content[i+1])
			if n2 > 9 || err != nil {
				panic("Invalid numeric mode content")
			}
			n3, err := alphanumericCode(content[i+2])
			if n3 > 9 || err != nil {
				panic("Invalid numeric mode content")
			}
			bits.Append(int(n1)*100+int(n2)*10+int(n3), 10)
		}
		switch len(content) % 3 {
		case 1:
			n1, err := alphanumericCode(content[len(content)-1])
			if n1 > 9 || err != nil {
				panic("Invalid numeric mode content")
			}
			bits.Append(int(n1), 4)
		case 2:
			n1, err := alphanumericCode(content[len(content)-2])
			if n1 > 9 || err != nil {
				panic("Invalid numeric mode content")
			}
			n2, err := alphanumericCode(content[len(content)-1])
			if n2 > 9 || err != nil {
				panic("Invalid numeric mode content")
			}
			bits.Append(int(n1)*10+int(n2), 7)
		}
	case modeAlphanumeric:
		for i := 0; i+1 < len(content); i += 2 {
			n1, err := alphanumericCode(content[i])
			if err != nil {
				panic("Invalid alphanumeric mode content")
			}
			n2, err := alphanumericCode(content[i+1])
			if err != nil {
				panic("Invalid alphanumeric mode content")
			}
			bits.Append(int(n1)*45+int(n2), 11)
		}
		if len(content)%2 == 1 {
			n1, err := alphanumericCode(content[len(content)-1])
			if err != nil {
				panic("Invalid alphanumeric mode content")
			}
			bits.Append(int(n1), 6)
		}
	case modeByte:
		for _, b := range content {
			bits.Append(int(b), 8)
		}
	default:
		panic("Unsupported mode")
	}
}

func appendTerminator(capacityBytes int, bits *BitVector) {
	capacity := capacityBytes * 8
	if bits.Length() > capacity {
		panic("bits.Length() > capacity")
	}
	for i := 0; i < 4 && bits.Length() < capacity; i++ {
		bits.AppendBit(false)
	}
	if bits.Length()%8 != 0 {
		for i := bits.Length() % 8; i < 8; i++ {
			bits.AppendBit(false)
		}
	}
	for {
		if bits.Length() >= capacity {
			break
		}
		bits.Append(0xec, 8)
		if bits.Length() >= capacity {
			break
		}
		bits.Append(0x11, 8)
	}
	if bits.Length() != capacity {
		panic("bits.Length() != capacity")
	}
}
