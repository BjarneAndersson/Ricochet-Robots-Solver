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

func HasNeighbor(currentNode byte, direction string) bool {
	switch direction {
	case "up":
		if bitOperations.HasBit(currentNode, 3) {
			return true
		}
	case "down":
		if bitOperations.HasBit(currentNode, 2) {
			return true
		}
	case "left":
		if bitOperations.HasBit(currentNode, 1) {
			return true
		}
	case "right":
		if bitOperations.HasBit(currentNode, 0) {
			return true
		}
	}
	return false
}

func GetNeighborNodePosition(position types.Position, direction string) types.Position {
	switch direction {
	case "up":
		position.Row -= 1
	case "down":
		position.Row += 1
	case "left":
		position.Column -= 1
	case "right":
		position.Column += 1
	}
	return position
}

func ConvTargetPositionToPosition(target uint16) (position types.Position) {
	targetBytePosition := target & 255
	return ConvBytePositionToPosition(byte(targetBytePosition))
}

func ConvBytePositionToPosition(bytePosition byte) (position types.Position) {
	position = types.Position{
		Column: (bytePosition & (15 << 4)) >> 4,
		Row:    bytePosition & 15,
	}
	return position
}

func GetMoveCount(b byte) byte {
	b = ((7 << 5) & b) >> 5
	return b
}
