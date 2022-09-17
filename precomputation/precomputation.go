package precomputation

import (
	"../helper"
	"../types"
)

func PrecomputeBoard(board *types.Board) (err error) {
	var status [16][16]bool

	targetPosition := helper.ConvTargetPositionToPosition(board.Target)
	pTargetNode := &(board.Board[targetPosition.Row][targetPosition.Column])
	setMoveCount(pTargetNode, 0)
	status[targetPosition.Row][targetPosition.Column] = true

	done := false

	for !done {
		done = true
		for indexRow, boardRow := range board.Board {
			for indexColumn := range boardRow {
				if !status[indexRow][indexColumn] {
					continue
				}

				node := board.Board[indexRow][indexColumn]

				status[indexRow][indexColumn] = false
				depth := helper.GetMoveCount(node) + 1

				preIndex := types.Position{Row: uint8(indexRow), Column: uint8(indexColumn)}

				for _, direction := range []string{"up", "down", "left", "right"} {
					index := preIndex

					for helper.HasNeighbor(board.Board[index.Row][index.Column], direction) {
						//pNeighborNode := getNeighborNode(board, &(board.Board[index.Row][index.Column]), direction)
						index = helper.GetNeighborNodePosition(index, direction)

						if helper.GetMoveCount(board.Board[index.Row][index.Column]) > depth {
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

func setMoveCount(pB *byte, moveCount uint8) {
	*pB = (31 & *pB) | (moveCount << 5)
}

func PrecomputeRobotMoves(board *types.Board) (robotMoves map[types.Position]map[string]types.Position, err error) {
	robotMoves = make(map[types.Position]map[string]types.Position)

	for rowIndex := 0; rowIndex < 16; rowIndex++ {
		for columnIndex := 0; columnIndex < 16; columnIndex++ {
			robotMoves[types.Position{Row: uint8(rowIndex), Column: uint8(columnIndex)}] = map[string]types.Position{"left": {}, "right": {}, "top": {}, "bottom": {}}
			for _, direction := range []string{"left", "right", "top", "bottom"} {
				robotMoves[types.Position{Row: uint8(rowIndex), Column: uint8(columnIndex)}][direction] = calculateRobotStoppingPosition(board, types.Position{Row: uint8(rowIndex), Column: uint8(columnIndex)}, direction)
			}
		}
	}

	return robotMoves, nil
}

func calculateRobotStoppingPosition(board *types.Board, startNodePosition types.Position, direction string) types.Position {
	cNodePosition := startNodePosition
	cNode := board.Board[cNodePosition.Row][cNodePosition.Column]

	for helper.HasNeighbor(cNode, direction) {
		cNode = board.Board[cNodePosition.Row][cNodePosition.Column]

		switch direction {
		case "left":
			cNodePosition = types.Position{Column: cNodePosition.Column - 1, Row: cNodePosition.Row}
		case "right":
			cNodePosition = types.Position{Column: cNodePosition.Column + 1, Row: cNodePosition.Row}
		case "top":
			cNodePosition = types.Position{Column: cNodePosition.Column, Row: cNodePosition.Row - 1}
		case "bottom":
			cNodePosition = types.Position{Column: cNodePosition.Column, Row: cNodePosition.Row + 1}
		}
	}

	return cNodePosition
}
