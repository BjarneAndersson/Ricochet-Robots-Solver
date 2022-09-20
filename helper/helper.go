package helper

import (
	"../bitOperations"
	"../types"
)

func GetTargetColor(target uint16) (color string, err error) {
	colorInBit := uint8((target & (15 << 12)) >> 12)

	if bitOperations.HasBit(colorInBit, 3) {
		return "yellow", nil
	} else if bitOperations.HasBit(colorInBit, 2) {
		return "red", nil
	} else if bitOperations.HasBit(colorInBit, 1) {
		return "green", nil
	} else if bitOperations.HasBit(colorInBit, 0) {
		return "blue", nil
	}
	return "", err
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

func SeparateRobots(boardState types.BoardState) (robots [4]byte) {
	robots[0] = uint8((uint32(boardState) & (uint32(255) << 24)) >> 24)
	robots[1] = uint8((uint32(boardState) & (uint32(255) << 16)) >> 16)
	robots[2] = uint8((uint32(boardState) & (uint32(255) << 8)) >> 8)
	robots[3] = uint8((uint32(boardState) & (uint32(255) << 0)) >> 0)
	return robots
}

func GetRobotColorByIndex(robotColors types.RobotColors, index uint8) (string, error) {
	for colorName, colorIndex := range robotColors {
		if index == colorIndex {
			return colorName, nil
		}
	}
	return "", fmt.Errorf("index %d not found", index)
}
