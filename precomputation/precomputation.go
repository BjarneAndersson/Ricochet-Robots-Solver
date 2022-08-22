package precomputation

import (
	"../helper"
	"../types"
	"fmt"
)

func PrecomputeBoard(board *types.Board) (err error) {
	var status [16][16]bool

	targetPosColumn, targetPosRow := convTargetPositionToPosition(board.Target)
	pTargetNode := &(board.Board[targetPosRow][targetPosColumn])
	setMoveCount(pTargetNode, 0)
	status[targetPosRow][targetPosColumn] = true

	done := false

	for !done {
		done = true
		for indexRow, boardRow := range board.Board {
			for indexColumn := range boardRow {
				if !status[indexRow][indexColumn] {
					continue
				}

				node := board.Board[indexRow][indexColumn]

				fmt.Printf("%v\n", node)

				status[indexRow][indexColumn] = false
				depth := getMoveCount(node) + 1

				preIndex := types.Position{Row: uint8(indexRow), Column: uint8(indexColumn)}

				for _, direction := range []string{"up", "down", "left", "right"} {
					index := preIndex

					for hasNeighbor(board.Board[index.Row][index.Column], direction) {
						//pNeighborNode := getNeighborNode(board, &(board.Board[index.Row][index.Column]), direction)
						index = getNeighborNodePosition(index, direction)

						fmt.Printf("%+v | %v - %v\n", index, getMoveCount(board.Board[index.Row][index.Column]), depth)

						if getMoveCount(board.Board[index.Row][index.Column]) > depth {
							setMoveCount(&(board.Board[index.Row][index.Column]), depth)
							status[index.Row][index.Column] = true
							done = false
						}
					}
				}
			}
		}
	}
	return nil
}

func hasNeighbor(currentNode byte, direction string) bool {
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

func getNeighborNodePosition(position types.Position, direction string) types.Position {
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

func convTargetPositionToPosition(target uint16) (column uint8, row uint8) {
	targetPosition := target & 255
	return convBytePositionToPosition(byte(targetPosition))
}

func convBytePositionToPosition(position byte) (column uint8, row uint8) {
	column = (position & (15 << 4)) >> 4
	row = position & 15
	return column, row
}

func getMoveCount(b byte) byte {
	b = ((7 << 5) & b) >> 5
	return b
}

func setMoveCount(pB *byte, moveCount uint8) {
	*pB = (31 & *pB) | (moveCount << 5)
}
