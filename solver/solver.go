package solver

import (
	"../helper"
	"../priorityQueue"
	"../types"
)

func reconstructPath(cameFrom map[types.BoardState]types.BoardState, currentBoardState types.BoardState) (path []types.BoardState, err error) {
	path = append(path, currentBoardState)

	for {
		if _, ok := cameFrom[currentBoardState]; ok {
			currentBoardState = cameFrom[currentBoardState]
			path = append(path, currentBoardState)
		} else {
			break
		}
	}

	// reverse the path
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path, nil
}

func Solver(board *types.Board, initBoardState types.BoardState) ([]types.BoardState, error) {
	openSet := make(priorityQueue.PriorityQueue, 1)

	openSet[0] = priorityQueue.Item{
		Value:    initBoardState,
		Priority: 0,
	}
	heap.Init(&openSet)

	cameFrom := make(map[types.BoardState]types.BoardState)

	var searchDepth uint8 // g score

	fScore := make(map[types.BoardState]uint8)
	fScore[initBoardState] = calcFScore(board, initBoardState, searchDepth)

	for openSet.Len() > 0 {
		currentBoardState := heap.Pop(&openSet).(*priorityQueue.Item)

		for indexRobot, robot := range helper.SeparateRobots(currentBoardState.Value) {
			robotPosition := helper.ConvBytePositionToPosition(robot)
			node := board.Board[robotPosition.Row][robotPosition.Column]

			for _, direction := range []string{"up", "down", "left", "right"} {
				cNode := node

				calculateStoppingPosition(board, cNode, direction)

				if cNode == node {
					// move robot
					newRobots := moveRobot(helper.SeparateRobots(currentBoardState.Value), uint8(indexRobot), direction)

					// create a new board state
					newBoardState := createNewBoardState(newRobots)

					// check if the new board state is the target
					// break -> reconstruct path
					if isRobotOnTarget(newBoardState, board.Target) {
						// add board state to cameFrom
						cameFrom[newBoardState] = currentBoardState.Value
						return reconstructPath(cameFrom, newBoardState)
					}

					// calc fScore for the new board state
					currentFScore := calcFScore(board, newBoardState, searchDepth)

					// check if the new board state is already in the queue
					isBoardStateInOpenSet(openSet, newBoardState)

					// add board state to cameFrom
					cameFrom[newBoardState] = currentBoardState.Value

					// add the new board state to the queue
					openSet.Push(
						priorityQueue.Item{
							Value:    newBoardState,
							Priority: int(currentFScore),
						})
				}

			}
		}
	}
	return []types.BoardState{}, nil
}

func calcFScore(board *types.Board, boardState types.BoardState, gScore uint8) (fScore uint8) {
	fScore = gScore + calcHScore(board, boardState)
	return fScore
}

func calcHScore(board *types.Board, boardState types.BoardState) (gScore uint8) {
	activeRobotPosition := helper.ConvBytePositionToPosition(uint8(boardState & (255 << 24)))

	node := board.Board[activeRobotPosition.Row][activeRobotPosition.Column]

	gScore = helper.GetMoveCount(node)
	return gScore
}

func createNewBoardState(robots [4]byte) types.BoardState {
	return types.BoardState(uint32(robots[0])<<24 | uint32(robots[1])<<16 | uint32(robots[2])<<8 | uint32(robots[3])<<0)
}

func calculateStoppingPosition(board *types.Board, startNode byte, direction string) types.Position {
	cNode := startNode
	for helper.HasNeighbor(cNode, direction) {
		cNodePosition := helper.ConvBytePositionToPosition(cNode)

		switch direction {
		case "left":
			cNode = board.Board[cNodePosition.Row][cNodePosition.Column-1]
		case "right":
			cNode = board.Board[cNodePosition.Row][cNodePosition.Column+1]
		case "up":
			cNode = board.Board[cNodePosition.Row-1][cNodePosition.Column]
		case "down":
			cNode = board.Board[cNodePosition.Row+1][cNodePosition.Column]
		}
	}
	return helper.ConvBytePositionToPosition(cNode)
}

func moveRobot(robots [4]byte, robotIndex uint8, direction string) (newRobots [4]byte) {
	newRobots = robots
	cRobot := &(newRobots[robotIndex])
	cRobotPosition := helper.ConvBytePositionToPosition(*cRobot)

	switch direction {
	case "left":
		cRobotPosition.Column -= 1
	case "right":
		cRobotPosition.Column += 1
	case "up":
		cRobotPosition.Row -= 1
	case "down":
		cRobotPosition.Row += 1
	}

	helper.ConvPosToByte(cRobot, cRobotPosition.Column, cRobotPosition.Row)
	return newRobots
}

func isBoardStateInOpenSet(openSet priorityQueue.PriorityQueue, boardState types.BoardState) bool {
	for _, iterateBoardState := range openSet {
		if iterateBoardState.Value == boardState {
			return true
		}
	}
	return false
}

func isRobotOnTarget(boardState types.BoardState, target uint16) bool {
	targetPosition := helper.ConvBytePositionToPosition(byte(target & 255))
	activeRobotPosition := helper.ConvBytePositionToPosition(byte(boardState & (255 << 24)))
	return activeRobotPosition.Column == targetPosition.Column && activeRobotPosition.Row == targetPosition.Row
}
