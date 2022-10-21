package helper

import (
	"Ricochet-Robot-Solver/internal/types"
	"fmt"
	"math"
)

// ConvPositionToByte Converts position to byte position: column (bit 7 to 4) | row (bit 3 to 0)
func ConvPositionToByte(position types.Position) byte {
	return position.Column<<4 + position.Row
}

// ConvBytePositionToPosition Returns the formatted struct position of the byte position
func ConvBytePositionToPosition(bytePosition byte) (position types.Position) {
	position = types.Position{
		Column: (bytePosition & (15 << 4)) >> 4,
		Row:    bytePosition & 15,
	}
	return position
}

// GetTargetColor Returns the color name of the target
func GetTargetColor(target types.Target) (color string, err error) {
	mask := uint16(math.Pow(2, 4)-1) << 12
	colorInBits := uint8((uint16(target) & mask) >> 12)

	if HasBit(colorInBits, 3) {
		return "yellow", nil
	} else if HasBit(colorInBits, 2) {
		return "red", nil
	} else if HasBit(colorInBits, 1) {
		return "green", nil
	} else if HasBit(colorInBits, 0) {
		return "blue", nil
	}
	return "", fmt.Errorf("target has no color")
}

// ConvTargetPositionToPosition Returns the position of the target
func ConvTargetPositionToPosition(target types.Target) (position types.Position) {
	targetBytePosition := target & 255
	return ConvBytePositionToPosition(byte(targetBytePosition))
}

// HasNeighbor Check if node has neighbor in the given direction
func HasNeighbor(currentNode types.Node, direction string) bool {
	switch direction {
	case "top":
		if HasBit(byte(currentNode), 3) {
			return true
		}
	case "bottom":
		if HasBit(byte(currentNode), 2) {
			return true
		}
	case "left":
		if HasBit(byte(currentNode), 1) {
			return true
		}
	case "right":
		if HasBit(byte(currentNode), 0) {
			return true
		}
	}
	return false
}

// GetNeighborNodePosition Returns the position of the neighbor node in the given direction
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

// SeparateRobots Separate robot into ordered array
func SeparateRobots(boardState types.BoardState) (robots [4]types.Robot) {
	robots[0] = types.Robot(uint8((uint32(boardState) & (uint32(255) << 24)) >> 24))
	robots[1] = types.Robot(uint8((uint32(boardState) & (uint32(255) << 16)) >> 16))
	robots[2] = types.Robot(uint8((uint32(boardState) & (uint32(255) << 8)) >> 8))
	robots[3] = types.Robot(uint8((uint32(boardState) & (uint32(255) << 0)) >> 0))
	return robots
}

// GetRobotColorCodeByIndex Returns the color code of the specified robot
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

// SetRobotColorByIndex Sets the robot color code by the index and the color name
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

// SetRobotColorCodeByIndex Sets the robot color code by the index
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

// GetMoveCount Get the minimum number of moves that a robot has to make in order to get to the target cell, if it could stop everywhere.
func GetMoveCount(b types.Node) uint8 {
	return uint8(((7 << 5) & b) >> 5)
}
