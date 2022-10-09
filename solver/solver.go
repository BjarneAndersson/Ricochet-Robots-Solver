package solver

import (
	"../config"
	"../helper"
	"../output"
	"../priorityQueue"
	"../tracker"
	"../types"
	"fmt"
	"sort"
)

func reconstructPath(cameFrom []uint64, endBoardState types.BoardState) (path []types.BoardState, err error) {
	path = append(path, endBoardState)
	currentBoardState := endBoardState

	for {
	startLoop:
		for _, currentPair := range cameFrom {
			if currentBoardState == types.BoardState((currentPair&(uint64(4294967295)<<32))>>32) {
				currentBoardState = types.BoardState(currentPair & 4294967295)
				path = append(path, currentBoardState)
				goto startLoop
			}
		}
		break
	}

	// reverse the path
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path, nil
}

func Solver(gameRound *types.GameRound, initBoardState types.BoardState, robotStoppingPositions *types.RobotStoppingPositions, conf config.Config) (tracker.TrackingDataSolver, []types.BoardState, error) {
	trackingData := tracker.TrackingDataSolver{}

	openSet := make(priorityQueue.PriorityQueue, 1)
	closedSet := make([]types.BoardState, 0)

	openSet[0] = priorityQueue.Item{
		Value:      initBoardState,
		HAndGScore: 0,
	}
	trackingData.InitializedBoardStates += 1

	cameFrom := make([]uint64, 0)

	for openSet.Len() > 0 {
		currentPriorityQueueItem := priorityQueue.Pop(&openSet)
		currentBoardState := currentPriorityQueueItem.Value

		if conf.Modes[conf.Mode]["output"].BoardStates == true {
			err := output.BoardState(currentBoardState, trackingData)
			if err != nil {
				return trackingData, []types.BoardState{}, err
			}
		}

		for indexRobot, robot := range helper.SeparateRobots(currentBoardState) {
			robotPosition := helper.ConvBytePositionToPosition(robot)
			node := gameRound.Board[robotPosition.Row][robotPosition.Column]
			nodePosition := types.Position{Column: robotPosition.Column, Row: robotPosition.Row}

			for _, direction := range []string{"top", "bottom", "left", "right"} {
				cNodePosition := calculateStoppingPosition(robotStoppingPositions, currentBoardState, nodePosition, direction)
				cNode := gameRound.Board[cNodePosition.Row][cNodePosition.Column]

				if cNode != node {
					// robot can be moved into direction

					// move robot
					newRobots := moveRobot(helper.SeparateRobots(currentBoardState), uint8(indexRobot), cNodePosition)

					// create a new gameRound state
					newBoardState := createNewBoardState(newRobots)

					// check if the new gameRound state is already in the queue
					if isBoardStateInOpenSet(&openSet, newBoardState) || isBoardStateInClosedSet(&closedSet, newBoardState) {
						continue
					}

					trackingData.InitializedBoardStates += 1

					gScoreNewBoardState := priorityQueue.GetGScore(currentPriorityQueueItem.HAndGScore) + 1
					hScoreNewBoardState := calcHScore(gameRound, newBoardState)

					// add gameRound state to cameFrom
					cameFrom = append(cameFrom, (uint64(newBoardState)<<32)|uint64(currentBoardState))

					// check if the new gameRound state is the target
					// break -> reconstruct path
					if hScoreNewBoardState == 0 {
						trackingData.EvaluatedBoardStates += 1
						path, err := reconstructPath(cameFrom, newBoardState)
						return trackingData, path, err
					}

					// add the new gameRound state to the queue
					openSet.Push(
						priorityQueue.Item{
							Value:      newBoardState,
							HAndGScore: priorityQueue.CombineHAndGScore(gScoreNewBoardState, hScoreNewBoardState),
						})
				}

			}
		}
		trackingData.EvaluatedBoardStates += 1
		closedSet = append(closedSet, currentBoardState)
	}
	return trackingData, []types.BoardState{}, fmt.Errorf("no route found")
}

func calcFScore(gameRound *types.GameRound, boardState types.BoardState, gScore uint8) (fScore uint8) {
	fScore = gScore + calcHScore(gameRound, boardState)
	return fScore
}

func calcHScore(gameRound *types.GameRound, boardState types.BoardState) (hScore uint8) {
	activeRobotPosition := helper.ConvBytePositionToPosition(uint8((boardState & (255 << 24)) >> 24))

	node := gameRound.Board[activeRobotPosition.Row][activeRobotPosition.Column]

	hScore = helper.GetMoveCount(node)
	return hScore
}

func createNewBoardState(robots [4]byte) types.BoardState {
	tRobots := robots

	var robotsSlice = tRobots[1:4]
	sort.Slice(robotsSlice, func(i, j int) bool {
		return robotsSlice[i] < robotsSlice[j]
	})
	return types.BoardState(uint32(robots[0])<<24 | uint32(robotsSlice[0])<<16 | uint32(robotsSlice[1])<<8 | uint32(robotsSlice[2])<<0)
}

func calculateStoppingPosition(robotStoppingPositions *types.RobotStoppingPositions, boardState types.BoardState, startNodePosition types.Position, direction string) (stoppingPosition types.Position) {
	switch direction {
	case "top":
		stoppingPosition = helper.ConvBytePositionToPosition(byte(((*robotStoppingPositions)[startNodePosition.Row][startNodePosition.Column] & (uint32(255) << 24)) >> 24))
	case "bottom":
		stoppingPosition = helper.ConvBytePositionToPosition(byte(((*robotStoppingPositions)[startNodePosition.Row][startNodePosition.Column] & (uint32(255) << 16)) >> 16))
	case "left":
		stoppingPosition = helper.ConvBytePositionToPosition(byte(((*robotStoppingPositions)[startNodePosition.Row][startNodePosition.Column] & (uint32(255) << 8)) >> 8))
	case "right":
		stoppingPosition = helper.ConvBytePositionToPosition(byte(((*robotStoppingPositions)[startNodePosition.Row][startNodePosition.Column] & (uint32(255) << 0)) >> 0))
	}

	if stoppingPosition == startNodePosition {
		return stoppingPosition
	}

	for _, robot := range helper.SeparateRobots(boardState) {
		robotPosition := helper.ConvBytePositionToPosition(robot)

		if robotPosition == startNodePosition {
			continue
		}

		switch direction {
		case "left":
			if robotPosition.Row == stoppingPosition.Row && robotPosition.Column >= stoppingPosition.Column && robotPosition.Column < startNodePosition.Column {
				stoppingPosition.Column = robotPosition.Column + 1
			}
		case "right":
			if robotPosition.Row == stoppingPosition.Row && robotPosition.Column <= stoppingPosition.Column && robotPosition.Column > startNodePosition.Column {
				stoppingPosition.Column = robotPosition.Column - 1
			}
		case "top":
			if robotPosition.Column == stoppingPosition.Column && robotPosition.Row >= stoppingPosition.Row && robotPosition.Row < startNodePosition.Row {
				stoppingPosition.Row = robotPosition.Row + 1
			}
		case "bottom":
			if robotPosition.Column == stoppingPosition.Column && robotPosition.Row <= stoppingPosition.Row && robotPosition.Row > startNodePosition.Row {
				stoppingPosition.Row = robotPosition.Row - 1
			}
		}
	}
	return stoppingPosition
}

func moveRobot(robots [4]byte, robotIndex uint8, endPosition types.Position) (newRobots [4]byte) {
	newRobots = robots

	newRobots[robotIndex] = helper.ConvPosToByte(endPosition)
	return newRobots
}

func isBoardStateInOpenSet(openSet *priorityQueue.PriorityQueue, boardState types.BoardState) bool {
	for _, iterateBoardState := range *openSet {
		if iterateBoardState.Value == boardState {
			return true
		}
	}
	return false
}

func isBoardStateInClosedSet(closedSet *[]types.BoardState, boardState types.BoardState) bool {
	for _, iterateBoardState := range *closedSet {
		if iterateBoardState == boardState {
			return true
		}
	}
	return false
}

func isRobotOnTarget(boardState *types.BoardState, target uint16) bool {
	targetPosition := helper.ConvBytePositionToPosition(byte(target & 255))
	activeRobotPosition := helper.ConvBytePositionToPosition(byte((*boardState & (255 << 24)) >> 24))
	return activeRobotPosition == targetPosition
}
