package bitOperations

// SetBit Sets the bit at pos in the byte pB.
func SetBit(pB *byte, position uint8, value bool) {
	if value {
		n := int(*pB)
		n |= boolToInteger(value) << position
		*pB = byte(n)
	} else {
		clearBit(pB, position)
	}
}

// Clears the bit at pos in n.
func clearBit(pB *uint8, position uint8) {
	mask := ^(1 << position)
	*pB &= uint8(mask)
}

// HasBit Checks whether the bit at pos n is set.
func HasBit(b byte, position uint8) bool {
	n := int(b)
	val := n & (1 << position)
	return val > 0
}

func boolToInteger(b bool) int {
	if b {
		return 1
	}
	return 0
}
