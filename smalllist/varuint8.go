package smalllist

const (
	maskSelectAllFirst7 = 0x7F
	maskSelect8th       = 0x80
)

func Encode(buffer []uint8, val int32) {
	var code int32
	i := 0
	for val != 0 {
		code = val & maskSelectAllFirst7
		if code != val {
			// setup continuation bit
			code |= maskSelect8th
		}
		buffer[i] = uint8(code)
		i += 1
		val = val >> 7
	}
}

func Decode(buffer []uint8, index int) int32 {
	var val int32
	hasContinuationBit := true
	i := index
	for hasContinuationBit {
		val += int32(buffer[i]&maskSelectAllFirst7) << ((i - index) * 7)
		if int(buffer[i])&maskSelect8th == 0 {
			hasContinuationBit = false
		}
		i += 1
	}
	return val
}
