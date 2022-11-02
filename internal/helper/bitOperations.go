package helper

// SetBit Sets the bit at the given position of the byte to the transferred value.
func SetBit(pB *byte, position uint8, value bool) {
	if value {
		n := int(*pB)
		n |= boolToInteger(value) << position
		*pB = byte(n)
	} else {
		clearBit(pB, position)
	}
}

// clearBit Sets the bit at the given position of the byte to 0.
func clearBit(pB *uint8, position uint8) {
	mask := ^(1 << position)
	*pB &= uint8(mask)
}

// HasBit Checks whether the bit at the position is set (1).
func HasBit(b byte, position uint8) bool {
	n := int(b)
	val := n & (1 << position)
	return val > 0
}

// boolToInteger Converts a boolean to an int
func boolToInteger(b bool) int {
	if b {
		return 1
	}
	return 0
}
