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

func ConvPosToByte(position types.Position) byte {
	return position.Column<<4 + position.Row
}

func HasNeighbor(currentNode byte, direction string) bool {
	switch direction {
	case "top":
		if bitOperations.HasBit(currentNode, 3) {
			return true
		}
	case "bottom":
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
	case "top":
		position.Row -= 1
	case "bottom":
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

func GetRobotColorCodeByIndex(robotOrder byte, index uint8) (robotColor types.RobotColor) {
	switch index {
	case 0:
		robotColor = types.RobotColor(robotOrder & (3 << 6) >> 6)
	case 1:
		robotColor = types.RobotColor(robotOrder & (3 << 4) >> 4)
	case 2:
		robotColor = types.RobotColor(robotOrder & (3 << 2) >> 2)
	case 3:
		robotColor = types.RobotColor(robotOrder & (3 << 0) >> 0)
	}
	return robotColor
}

func SetRobotColorByIndex(robotOrder *byte, colorName string, index uint8) {
	var colorCode types.RobotColor

	switch colorName {
	case "yellow":
		colorCode = types.RobotColorYellow
	case "red":
		colorCode = types.RobotColorRed
	case "green":
		colorCode = types.RobotColorGreen
	case "blue":
		colorCode = types.RobotColorBlue
	}

	SetRobotColorCodeByIndex(robotOrder, colorCode, index)
}

func SetRobotColorCodeByIndex(robotOrder *byte, colorCode types.RobotColor, index uint8) {
	switch index {
	case 0:
		*robotOrder = (*robotOrder) | (byte(colorCode) << 6)
	case 1:
		*robotOrder = (*robotOrder) | (byte(colorCode) << 4)
	case 2:
		*robotOrder = (*robotOrder) | (byte(colorCode) << 2)
	case 3:
		*robotOrder = (*robotOrder) | (byte(colorCode) << 0)
	}
}
