package bitOperations

// SetBit Sets the bit at pos in the byte pB.
func SetBit(pB *byte, position uint8, value bool) {
	n := int(*pB)
	n |= boolToInteger(value) << position
	*pB = byte(n)
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
