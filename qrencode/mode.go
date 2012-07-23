package qrencode

const (
	modeTerminator         = modeIndicator(0)
	modeNumeric            = modeIndicator(1)
	modeAlphanumeric       = modeIndicator(2)
	modeStructuredAppend   = modeIndicator(3)
	modeByte               = modeIndicator(4)
	modeFNC1FirstPosition  = modeIndicator(5)
	modeECI                = modeIndicator(7)
	modeKanji              = modeIndicator(8)
	modeFNC1SecondPosition = modeIndicator(9)
)

type modeIndicator int

func getMode(content string) modeIndicator {
	numeric := false
	alphanumeric := false
	for _, b := range []byte(content) {
		if code, err := alphanumericCode(b); err != nil {
			return modeByte
		} else if code < 10 {
			numeric = true
		} else {
			alphanumeric = true
		}
	}
	if alphanumeric {
		return modeAlphanumeric
	}
	if numeric {
		return modeNumeric
	}
	return modeByte
}

func (m modeIndicator) characterCountBits(v versionNumber) int {
	switch {
	case v < 1:
		panic("Invalid versionNumber")
	case v <= 9:
		switch m {
		case modeNumeric:
			return 10
		case modeAlphanumeric:
			return 9
		case modeByte:
			return 8
		case modeKanji:
			return 8
		default:
			panic("Unsupported mode")
		}
	case v <= 26:
		switch m {
		case modeNumeric:
			return 12
		case modeAlphanumeric:
			return 11
		case modeByte:
			return 16
		case modeKanji:
			return 10
		default:
			panic("Unsupported mode")
		}
	case v <= 40:
		switch m {
		case modeNumeric:
			return 14
		case modeAlphanumeric:
			return 13
		case modeByte:
			return 16
		case modeKanji:
			return 12
		default:
			panic("Unsupported mode")
		}
	}
	panic("Invalid versionNumber")
}
