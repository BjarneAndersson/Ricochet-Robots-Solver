package output

import (
	"Ricochet-Robot-Solver/internal/helper"
	"Ricochet-Robot-Solver/internal/input"
	"Ricochet-Robot-Solver/internal/log"
	"Ricochet-Robot-Solver/internal/tracker"
	"Ricochet-Robot-Solver/internal/types"
	"encoding/json"
	"fmt"
)

// Neighbors Prints per node the direction of the neighbors
func Neighbors(board *types.Board, recorder *log.SolveRecorder) (err error) {
	var outputNodesAndNeighbors []OutNeighbors

	// convert the board to a json representation
	for rowIndex := range board.Grid {
		for columnIndex, node := range board.Grid[rowIndex] {
			outNode := OutNeighbors{}

			outNode.Node = types.Position{
				Column: byte(columnIndex),
				Row:    byte(rowIndex),
			}

			if !helper.HasNeighbor(node, "top") {
				outNode.Directions.Top = true
			}
			if !helper.HasNeighbor(node, "bottom") {
				outNode.Directions.Bottom = true
			}
			if !helper.HasNeighbor(node, "left") {
				outNode.Directions.Left = true
			}
			if !helper.HasNeighbor(node, "right") {
				outNode.Directions.Right = true
			}

			outputNodesAndNeighbors = append(outputNodesAndNeighbors, outNode)
		}
	}

	outputNodesAndNeighborsJson, err := json.Marshal(outputNodesAndNeighbors)
	if err != nil {
		fmt.Println("Error marshalling board to JSON:", err)
		return err
	}

	recorder.Record(fmt.Sprintf("NEIGHBORS: %v", string(outputNodesAndNeighborsJson)))
	return nil
}

// BoardState Prints the board state and current tracking data
func BoardState(boardState types.BoardState, trackingData tracker.TrackingDataSolver, recorder *log.SolveRecorder) (err error) {
	boardStateInBits, err := convertNumberToBits(int(boardState), 32)
	if err != nil {
		return err
	}

	recorder.Record(fmt.Sprintf("BOARD STATE: %v | %v\n", boardState, boardStateInBits))

	robots := helper.SeparateRobots(boardState)
	for _, robotPosition := range robots {
		recorder.Record(fmt.Sprintf("%+v\n", helper.ConvBytePositionToPosition(byte(robotPosition))))
	}
	recorder.Record(fmt.Sprintf("INITIALIZED STATES: %v | EVALUATED STATES: %v\n", trackingData.InitializedBoardStates, trackingData.EvaluatedBoardStates))
	return nil
}

// Path Prints out the path from start state to end state
func Path(path []types.BoardState, trackingData tracker.TrackingDataSolver, robotOrder types.RobotOrder, recorder *log.SolveRecorder) (err error) {
	recorder.Record(fmt.Sprintf("STATISTICS: Moves: %d | Time: %s | Initialized states: %v | Evaluated states: %v | States per second: %v",
		len(path)-1,
		trackingData.Duration,
		trackingData.InitializedBoardStates,
		trackingData.EvaluatedBoardStates,
		1000000000*trackingData.EvaluatedBoardStates/uint(trackingData.Duration)))
	recorder.Record("PATH")

	var outputPath []PathStep

	var previousBoardState types.BoardState
	var previousRobotOrder = robotOrder
	for indexBoardState, boardState := range path {
		robots := helper.SeparateRobots(boardState)
		outputRobots, err := convRobotsFromSolverToOutputRepresentation(robots, robotOrder)
		if err != nil {
			return err
		}

		cRobotOrder, err := getNewRobotOrder(previousBoardState, previousRobotOrder, boardState)
		if err != nil {
			return err
		}

		if indexBoardState == 0 { // initial position
			outputStringInitialPosition := "INITIAL POSITION: "

			previousBoardState = boardState

			for indexRobot, robotPosition := range robots {
				robotColor, err := GetRobotColorByIndex(robotOrder, uint8(indexRobot))
				if err != nil {
					return err
				}
				msg := fmt.Sprintf("%s: %+v ", robotColor, helper.ConvBytePositionToPosition(byte(robotPosition)))
				outputStringInitialPosition += msg
			}

			previousRobotOrder = robotOrder
			recorder.Record(outputStringInitialPosition)

			// add data to json output
			outputPath = append(outputPath, PathStep{
				Robots:         outputRobots,
				MovedRobot:     "",
				MovedDirection: "",
			})
		} else { // moved position / step
			robotColor, direction, err := getMovedRobotColorAndDirection(previousBoardState, boardState, cRobotOrder)
			if err != nil {
				return err
			}

			outputStringMove := fmt.Sprintf("Move: %v | ", indexBoardState)

			switch robotColor {
			case "yellow":
				outputStringMove += "Yellow: "
			case "red":
				outputStringMove += "Red: "
			case "green":
				outputStringMove += "Green: "
			case "blue":
				outputStringMove += "Blue: "
			}

			outputStringMove += fmt.Sprintf("%v", direction)

			recorder.Record(outputStringMove)

			previousBoardState = boardState
			previousRobotOrder = cRobotOrder

			// add data to json output
			outputPath = append(outputPath, PathStep{
				Robots:         outputRobots,
				MovedRobot:     robotColor,
				MovedDirection: direction,
			})
		}

		if indexBoardState == (len(path) - 1) { // end position
			outputStringFinalPosition := "FINAL POSITION: "

			for _, robotPosition := range robots {
				msg := fmt.Sprintf("%+v ", helper.ConvBytePositionToPosition(byte(robotPosition)))

				outputStringFinalPosition += msg
			}
			recorder.Record(outputStringFinalPosition)

			// add data to json output
			outputPath = append(outputPath, PathStep{
				Robots:         outputRobots,
				MovedRobot:     "",
				MovedDirection: "",
			})
		}
	}

	outputPathJson, err := json.Marshal(outputPath)
	if err != nil {
		fmt.Println("Error marshalling board to JSON:", err)
		return err
	}

	recorder.Record(fmt.Sprintf("PATH: %v", string(outputPathJson)))

	return nil
}

// RobotStoppingPositions Prints per node the stopping positions of a theoretical moved robot per direction
func RobotStoppingPositions(robotStoppingPositions *types.RobotStoppingPositions, recorder *log.SolveRecorder) (err error) {
	var outputRobotStoppingPositions OutRobotStoppingPositions

	for iRow, robotStoppingPositionsRow := range *robotStoppingPositions {
		for iColumn, stoppingPositions := range robotStoppingPositionsRow {
			outputRobotStoppingPositions = append(outputRobotStoppingPositions, OutRobotStoppingPosition{
				Node: types.Position{
					Column: byte(iColumn),
					Row:    byte(iRow),
				},
				Direction: OutDirectionAndPosition{
					Up:    helper.ConvBytePositionToPosition(byte(stoppingPositions >> 24)),
					Down:  helper.ConvBytePositionToPosition(byte(stoppingPositions >> 16)),
					Left:  helper.ConvBytePositionToPosition(byte(stoppingPositions >> 8)),
					Right: helper.ConvBytePositionToPosition(byte(stoppingPositions)),
				},
			})
		}
	}

	robotStoppingPositionsJson, err := json.Marshal(outputRobotStoppingPositions)
	if err != nil {
		fmt.Println("Error marshalling board to JSON:", err)
		return err
	}

	recorder.Record(fmt.Sprintf("ROBOT STOPPING POSITIONS: %v", string(robotStoppingPositionsJson)))

	return nil
}

// convertNumberToBits Converts a number into a representation as bits
func convertNumberToBits(number int, fill int) (string, error) {
	switch fill {
	case 4:
		return fmt.Sprintf("%04b", number), nil
	case 8:
		return fmt.Sprintf("%08b", number), nil
	case 32:
		return fmt.Sprintf("%32b", number), nil
	}
	return "", fmt.Errorf("invalid operation")
}

// GetRobotColorByIndex Returns the color name of the robot at the given index
func GetRobotColorByIndex(robotOrder types.RobotOrder, index uint8) (string, error) {
	robotColorCode := helper.GetRobotColorCodeByIndex(robotOrder, index)

	switch robotColorCode {
	case 0:
		return "yellow", nil
	case 1:
		return "red", nil
	case 2:
		return "green", nil
	case 3:
		return "blue", nil
	}
	return "", fmt.Errorf("index out of range")
}

// getMovedRobotColorAndDirection Returns the moved robot color name and the direction
func getMovedRobotColorAndDirection(previousBoardState types.BoardState, currentBoardState types.BoardState, currentRobotOrder types.RobotOrder) (robotColor string, direction string, err error) {
	preRobots := helper.SeparateRobots(previousBoardState)
	preRobotsNotContainedInCurRobots := preRobots[:]
	curRobots := helper.SeparateRobots(currentBoardState)

	var robotIndex uint8
	robotContainedInBothBoardStates := false

	for indexCurRobot, curRobot := range curRobots {
		robotContainedInBothBoardStates = false
		for _, preRobot := range preRobots {
			if curRobot == preRobot {
				robotContainedInBothBoardStates = true
				err := removeElement(&preRobotsNotContainedInCurRobots, preRobot)
				if err != nil {
					return "", "", err
				}
				break
			}
		}
		if robotContainedInBothBoardStates == false {
			robotIndex = uint8(indexCurRobot)
			break
		}
	}

	robotColor, err = GetRobotColorByIndex(currentRobotOrder, robotIndex)
	if err != nil {
		return "", "", err
	}

	direction, err = evaluateDirectionChange(preRobotsNotContainedInCurRobots[0], curRobots[robotIndex])
	if err != nil {
		return "", "", err
	}

	return robotColor, direction, nil
}

func removeElement(s *[]types.Robot, e types.Robot) error {
	for iterIndex, iterElement := range *s {
		if iterElement == e {
			*s = append((*s)[:iterIndex], (*s)[iterIndex+1:]...)
			return nil
		}
	}
	return fmt.Errorf("element: %v not in slice", e)
}

// evaluateDirectionChange Returns the direction the robot moved
func evaluateDirectionChange(previousRobot types.Robot, currentRobot types.Robot) (direction string, err error) {
	preRobotPos := helper.ConvBytePositionToPosition(byte(previousRobot))
	curRobotPos := helper.ConvBytePositionToPosition(byte(currentRobot))

	if preRobotPos.Column > curRobotPos.Column {
		return "left", nil
	} else if preRobotPos.Column < curRobotPos.Column {
		return "right", nil
	} else if preRobotPos.Row > curRobotPos.Row {
		return "up", nil
	} else if preRobotPos.Row < curRobotPos.Row {
		return "down", nil
	}

	return "", fmt.Errorf("robot position has not changed")
}

// getNewRobotOrder Calculated the new robot order based on the old order and the new robots
func getNewRobotOrder(previousBoardState types.BoardState, previousRobotOrder types.RobotOrder, currentBoardState types.BoardState) (currentRobotOrder types.RobotOrder, err error) {
	previousRobots := helper.SeparateRobots(previousBoardState)
	currentRobots := helper.SeparateRobots(currentBoardState)

	usedColorCodes := map[types.RobotColor]bool{0: false, 1: false, 2: false, 3: false}

	var indexCurrentNotExistingRobot int

	for indexCurrentRobot, currentRobot := range currentRobots {
		exist, indexPreviousRobotOrder := robotInRobots(currentRobot, previousRobots)

		if exist == false {
			indexCurrentNotExistingRobot = indexCurrentRobot
		} else {
			usedColorCodes[helper.GetRobotColorCodeByIndex(previousRobotOrder, uint8(indexPreviousRobotOrder))] = true
			helper.SetRobotColorCodeByIndex(&currentRobotOrder, helper.GetRobotColorCodeByIndex(previousRobotOrder, uint8(indexPreviousRobotOrder)), uint8(indexCurrentRobot))
		}
	}

	for colorCode, exist := range usedColorCodes {
		if exist == false {
			helper.SetRobotColorCodeByIndex(&currentRobotOrder, colorCode, uint8(indexCurrentNotExistingRobot))
		}
	}

	return currentRobotOrder, nil
}

// robotInRobots Checks if the robot is in the robots
func robotInRobots(robot types.Robot, robots [4]types.Robot) (exist bool, index int) {
	for iterIndexRobot, iterRobot := range robots {
		if iterRobot == robot {
			return true, iterIndexRobot
		}
	}
	return false, -1
}

func convRobotsFromSolverToOutputRepresentation(robots [4]types.Robot, robotOrder types.RobotOrder) (outputRobots input.Robots, err error) {
	for indexRobot, robotPosition := range robots {
		robotColor, err := GetRobotColorByIndex(robotOrder, uint8(indexRobot))
		if err != nil {
			return outputRobots, err
		}

		outputRobots[indexRobot] = input.Robot{
			Color:    robotColor,
			Position: helper.ConvBytePositionToPosition(byte(robotPosition)),
		}
	}

	return outputRobots, nil
}
