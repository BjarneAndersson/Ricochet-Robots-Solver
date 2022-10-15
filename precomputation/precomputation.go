package precomputation

import (
	"../helper"
	"../types"
)

func PrecomputeBoard(board *types.Board) (err error) {
	var status [16][16]bool

	targetPosition := helper.ConvTargetPositionToPosition(board.Target)
	pTargetNode := &(board.Grid[targetPosition.Row][targetPosition.Column])
	setMoveCount(pTargetNode, 0)
	status[targetPosition.Row][targetPosition.Column] = true

	done := false

	for !done {
		done = true
		for indexRow, boardRow := range board.Grid {
			for indexColumn := range boardRow {
				if !status[indexRow][indexColumn] {
					continue
				}

				node := board.Grid[indexRow][indexColumn]

				status[indexRow][indexColumn] = false
				depth := helper.GetMoveCount(node) + 1

				preIndex := types.Position{Row: uint8(indexRow), Column: uint8(indexColumn)}

				for _, direction := range []string{"top", "bottom", "left", "right"} {
					index := preIndex

					for helper.HasNeighbor(board.Grid[index.Row][index.Column], direction) {
						//pNeighborNode := getNeighborNode(board, &(board.Board[index.Row][index.Column]), direction)
						index = helper.GetNeighborNodePosition(index, direction)

						if helper.GetMoveCount(board.Grid[index.Row][index.Column]) > depth {
							setMoveCount(&(board.Grid[index.Row][index.Column]), depth)
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

func setMoveCount(pB *byte, moveCount uint8) {
	*pB = (31 & *pB) | (moveCount << 5)
}

func PrecomputeRobotMoves(board *types.Board) (robotStoppingPositions types.RobotStoppingPositions, err error) {
	for rowIndex := 0; rowIndex < 16; rowIndex++ {
		for columnIndex := 0; columnIndex < 16; columnIndex++ {
			for _, direction := range []string{"top", "bottom", "left", "right"} {
				switch direction {
				case "top":
					robotStoppingPositions[rowIndex][columnIndex] = robotStoppingPositions[rowIndex][columnIndex] | (uint32(calculateRobotStoppingPosition(board, types.Position{Row: uint8(rowIndex), Column: uint8(columnIndex)}, direction)) << 24)
				case "bottom":
					robotStoppingPositions[rowIndex][columnIndex] = robotStoppingPositions[rowIndex][columnIndex] | (uint32(calculateRobotStoppingPosition(board, types.Position{Row: uint8(rowIndex), Column: uint8(columnIndex)}, direction)) << 16)
				case "left":
					robotStoppingPositions[rowIndex][columnIndex] = robotStoppingPositions[rowIndex][columnIndex] | (uint32(calculateRobotStoppingPosition(board, types.Position{Row: uint8(rowIndex), Column: uint8(columnIndex)}, direction)) << 8)
				case "right":
					robotStoppingPositions[rowIndex][columnIndex] = robotStoppingPositions[rowIndex][columnIndex] | (uint32(calculateRobotStoppingPosition(board, types.Position{Row: uint8(rowIndex), Column: uint8(columnIndex)}, direction)) << 0)
				}
			}
		}
	}

	return robotStoppingPositions, nil
}

func calculateRobotStoppingPosition(board *types.Board, startNodePosition types.Position, direction string) byte {
	cNodePosition := startNodePosition
	cNode := board.Grid[cNodePosition.Row][cNodePosition.Column]

	for helper.HasNeighbor(cNode, direction) {

		switch direction {
		case "left":
			cNodePosition.Column -= 1
		case "right":
			cNodePosition.Column += 1
		case "top":
			cNodePosition.Row -= 1
		case "bottom":
			cNodePosition.Row += 1
		}
		cNode = board.Grid[cNodePosition.Row][cNodePosition.Column]
	}

	return helper.ConvPosToByte(cNodePosition)
}
