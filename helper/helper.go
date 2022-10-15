package helper

import (
	"../bitOperations"
	"../types"
	"fmt"
	"math"
)

func GetTargetColor(target uint16) (color string, err error) {
	mask := uint16(math.Pow(2, 4)-1) << 12
	colorInBits := uint8((target & mask) >> 12)

	if bitOperations.HasBit(colorInBits, 3) {
		return "yellow", nil
	} else if bitOperations.HasBit(colorInBits, 2) {
		return "red", nil
	} else if bitOperations.HasBit(colorInBits, 1) {
		return "green", nil
	} else if bitOperations.HasBit(colorInBits, 0) {
		return "blue", nil
	}
	return "", fmt.Errorf("target has no color")
}

// ConvPosToByte Converts position to byte position: column (bit 7 to 4) | row (bit 3 to 0)
func ConvPosToByte(position types.Position) byte {
	return position.Column<<4 + position.Row
}

// HasNeighbor Check if node has neighbor in the given direction
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

// GetMoveCount Get the minimum number of moves that a robot has to make in order to get to the target cell, if it could stop everywhere.
func GetMoveCount(b byte) byte {
	b = ((7 << 5) & b) >> 5
	return b
}

// SeparateRobots Separate robot into ordered array
func SeparateRobots(boardState types.BoardState) (robots [4]byte) {
	robots[0] = uint8((uint32(boardState) & (uint32(255) << 24)) >> 24)
	robots[1] = uint8((uint32(boardState) & (uint32(255) << 16)) >> 16)
	robots[2] = uint8((uint32(boardState) & (uint32(255) << 8)) >> 8)
	robots[3] = uint8((uint32(boardState) & (uint32(255) << 0)) >> 0)
	return robots
}

func GetRobotColorCodeByIndex(robotOrder types.RobotOrder, index uint8) (robotColor types.RobotColor) {
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

func SetRobotColorByIndex(robotOrder *types.RobotOrder, colorName string, index uint8) {
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

func SetRobotColorCodeByIndex(robotOrder *types.RobotOrder, colorCode types.RobotColor, index uint8) {
	switch index {
	case 0:
		*robotOrder = types.RobotOrder(byte(*robotOrder) | (byte(colorCode) << 6))
	case 1:
		*robotOrder = types.RobotOrder(byte(*robotOrder) | (byte(colorCode) << 4))
	case 2:
		*robotOrder = types.RobotOrder(byte(*robotOrder) | (byte(colorCode) << 2))
	case 3:
		*robotOrder = types.RobotOrder(byte(*robotOrder) | (byte(colorCode) << 0))
	}
}
