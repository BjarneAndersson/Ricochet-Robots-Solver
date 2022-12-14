package solver

import (
	"Ricochet-Robot-Solver/internal/config"
	"Ricochet-Robot-Solver/internal/helper"
	"Ricochet-Robot-Solver/internal/output"
	"Ricochet-Robot-Solver/internal/tracker"
	"Ricochet-Robot-Solver/internal/types"
	"container/heap"
	"fmt"
	"math"
	"sort"
)

// reconstructPath Reconstructs the path of board states, if a solution was found
func reconstructPath(cameFrom []uint64, endBoardState types.BoardState) (path []types.BoardState, err error) {
	path = append(path, endBoardState)
	currentBoardState := endBoardState

	maskToBoardState := uint64(math.Pow(2, 32) - 1)
	maskFromBoardState := uint64(math.Pow(2, 32)-1) << 32

	for {
	startLoop:
		for _, currentPair := range cameFrom {
			// if current board state is from-board-state
			if currentBoardState == types.BoardState((currentPair&maskFromBoardState)>>32) {
				// declare new current board state as the to-board-state
				currentBoardState = types.BoardState(currentPair & maskToBoardState)
				path = append(path, currentBoardState)
				goto startLoop
			}
		}
		// if path retracing is done
		break
	}

	// reverse the path
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path, nil
}

// Solver Solve the board
func Solver(board *types.Board, initBoardState types.BoardState, robotStoppingPositions *types.RobotStoppingPositions, conf config.Config) (tracker.TrackingDataSolver, []types.BoardState, error) {
	// initialization
	trackingData := tracker.TrackingDataSolver{}

	openSet := make(priorityQueue, 1)
	openSetMap := make(map[types.BoardState]struct{})
	closedSetMap := make(map[types.BoardState]struct{})

	cameFrom := make([]uint64, 0)

	// add initial board state to open set
	openSet[0] = &item{
		Value:      initBoardState,
		HAndGScore: 0,
		index:      0,
	}
	heap.Init(&openSet)
	trackingData.InitializedBoardStates += 1
	openSetMap[initBoardState] = struct{}{}

	for openSet.Len() > 0 {
		// get the item with the lowest f score
		currentPriorityQueueItem := heap.Pop(&openSet).(*item)
		currentBoardState := currentPriorityQueueItem.Value

		// output current board state based on the configuration
		if conf.Modes[conf.Mode]["output"].BoardStates == true {
			err := output.BoardState(currentBoardState, trackingData)
			if err != nil {
				return trackingData, []types.BoardState{}, err
			}
		}

		for indexRobot, robot := range helper.SeparateRobots(currentBoardState) {
			robotPosition := helper.ConvBytePositionToPosition(byte(robot))
			node := board.Grid[robotPosition.Row][robotPosition.Column]

			for _, direction := range []string{"top", "bottom", "left", "right"} {
				// get the stopping position of the robot in the given direction
				cNodePosition := calculateStoppingPosition(robotStoppingPositions, currentBoardState, robotPosition, direction)
				cNode := board.Grid[cNodePosition.Row][cNodePosition.Column]

				// if robot has been moved
				if cNode != node {
					newRobots := moveRobot(helper.SeparateRobots(currentBoardState), uint8(indexRobot), cNodePosition)

					newBoardState := CreateNewBoardState(newRobots)

					// check if the new board state is already in the queue or completely evaluated
					if isBoardStateInOpenSet(&openSetMap, newBoardState) || isBoardStateInClosedSet(&closedSetMap, newBoardState) {
						continue
					}

					trackingData.InitializedBoardStates += 1

					// add board state to cameFrom
					cameFrom = append(cameFrom, (uint64(newBoardState)<<32)|uint64(currentBoardState))

					// calculate g score: current g score + 1
					gScoreNewBoardState := getGScore(currentPriorityQueueItem.HAndGScore) + 1

					// calculate h score: prediction of minimal moves to go
					hScoreNewBoardState := calcHScore(board, newBoardState)

					// check if the active robot is on the target
					if hScoreNewBoardState == 0 {
						trackingData.EvaluatedBoardStates += 1
						// reconstruct the path
						path, err := reconstructPath(cameFrom, newBoardState)
						return trackingData, path, err
					}

					// add the new board state to the queue
					heap.Push(&openSet, &item{
						Value:      newBoardState,
						HAndGScore: combineHAndGScore(gScoreNewBoardState, hScoreNewBoardState),
					})
					openSetMap[newBoardState] = struct{}{}
				}

			}
		}
		trackingData.EvaluatedBoardStates += 1
		closedSetMap[currentBoardState] = struct{}{}
	}
	return trackingData, []types.BoardState{}, fmt.Errorf("no route found")
}

// calcHScore Gets the minimal move count of from the node the active robot is currently on
func calcHScore(board *types.Board, boardState types.BoardState) (hScore uint8) {
	activeRobotPosition := helper.ConvBytePositionToPosition(uint8((boardState & (255 << 24)) >> 24))

	// get node of active robot
	node := board.Grid[activeRobotPosition.Row][activeRobotPosition.Column]

	// get the minimal move count of that node
	hScore = helper.GetMoveCount(node)
	return hScore
}

// CreateNewBoardState Creates a new board state based on the robot positions
func CreateNewBoardState(robots [4]types.Robot) types.BoardState {
	tRobots := robots

	// sort none active robots
	var robotsSlice = tRobots[1:4]
	sort.Slice(robotsSlice, func(i, j int) bool {
		return robotsSlice[i] < robotsSlice[j]
	})
	// compose new board state
	return types.BoardState(uint32(robots[0])<<24 | uint32(robotsSlice[0])<<16 | uint32(robotsSlice[1])<<8 | uint32(robotsSlice[2])<<0)
}

// calculateStoppingPosition Calculates the stopping position from the given start position and direction
func calculateStoppingPosition(robotStoppingPositions *types.RobotStoppingPositions, boardState types.BoardState, startNodePosition types.Position, direction string) (stoppingPosition types.Position) {
	// get the precomputed stopping position of the given position and direction
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

	// gard condition: if robot have not been moved
	if stoppingPosition == startNodePosition {
		return stoppingPosition
	}

	// check if robots are in the path
	for _, robot := range helper.SeparateRobots(boardState) {
		robotPosition := helper.ConvBytePositionToPosition(byte(robot))

		if robotPosition == startNodePosition {
			continue
		}

		// if the robot is in the way -> move stopping position in front of the robot
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

// moveRobot Moves the robot to the end position
func moveRobot(robots [4]types.Robot, robotIndex uint8, endPosition types.Position) [4]types.Robot {
	robots[robotIndex] = types.Robot(helper.ConvPositionToByte(endPosition))
	return robots
}

// isBoardStateInOpenSet Checks if the board state is in the open set
func isBoardStateInOpenSet(openSet *map[types.BoardState]struct{}, boardState types.BoardState) bool {
	_, ok := (*openSet)[boardState]
	return ok
}

// isBoardStateInClosedSet Checks if the board state is in the closed set
func isBoardStateInClosedSet(closedSet *map[types.BoardState]struct{}, boardState types.BoardState) bool {
	_, ok := (*closedSet)[boardState]
	return ok
}
