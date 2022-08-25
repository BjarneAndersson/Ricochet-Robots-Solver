package helper

import (
	"../bitOperations"
	"../types"
)

func GetTargetColor(target uint16) (color types.Colors, err error) {
	colorInBit := uint8((target & (15 << 12)) >> 12)

	if bitOperations.HasBit(colorInBit, 3) {
		return types.Yellow, nil
	} else if bitOperations.HasBit(colorInBit, 2) {
		return types.Red, nil
	} else if bitOperations.HasBit(colorInBit, 1) {
		return types.Green, nil
	} else if bitOperations.HasBit(colorInBit, 0) {
		return types.Blue, nil
	}
	return 0, err
}

func ConvPosToByte(pB *byte, column uint8, row uint8) {
	*pB = column<<4 + row
}
