package output

import (
	"Ricochet-Robot-Solver/internal/helper"
	"Ricochet-Robot-Solver/internal/tracker"
	"Ricochet-Robot-Solver/internal/types"
	"fmt"
	"github.com/fatih/color"
)

// Neighbors Prints per node the direction of the neighbors
func Neighbors(board *types.Board) (err error) {
	fmt.Printf("\n\n====================\n")
	for rowIndex := range board.Grid {
		for _, node := range board.Grid[rowIndex] {
			outputString := fmt.Sprintf("")
			if !helper.HasNeighbor(node, "top") {
				outputString += "N"
			}
			if !helper.HasNeighbor(node, "bottom") {
				outputString += "S"
			}
			if !helper.HasNeighbor(node, "left") {
				outputString += "W"
			}
			if !helper.HasNeighbor(node, "right") {
				outputString += "E"
			}
			if len(outputString) == 0 {
				outputString = "A"
			}
			//pNode, err := convertNumberToBits(int(node&15), 4)
			if err != nil {
				return err
			}
			fmt.Printf("%v ", outputString)
		}
		fmt.Println("")
	}
	fmt.Printf("====================\n")
	return nil
}

// BoardState Prints the board state and current tracking data
func BoardState(boardState types.BoardState, trackingData tracker.TrackingDataSolver) (err error) {
	boardStateInBits, err := convertNumberToBits(int(boardState), 32)
	if err != nil {
		return err
	}

	fmt.Printf("\n\n====================\n")
	fmt.Printf("Board State: %v | %v\n", boardState, boardStateInBits)

	robots := helper.SeparateRobots(boardState)
	for _, robotPosition := range robots {
		fmt.Printf("%+v\n", helper.ConvBytePositionToPosition(byte(robotPosition)))
	}
	fmt.Printf("Initialized states: %v | Evaluated states: %v\n", trackingData.InitializedBoardStates, trackingData.EvaluatedBoardStates)
	fmt.Printf("====================\n")
	return nil
}

// Path Prints out the path from start state to end state
func Path(path []types.BoardState, trackingData tracker.TrackingDataSolver, robotOrder types.RobotOrder) (err error) {
	fmt.Printf("\n\n====================\n")
	fmt.Printf("Moves\t%d\n", len(path)-1)
	fmt.Printf("Time\t%s\n", trackingData.Duration)
	robotOrderInBits, _ := convertNumberToBits(int(robotOrder), 8)
	fmt.Printf("Path\t%+v\t%+v\n", path, robotOrderInBits)
	fmt.Printf("Initialized states: %v | Evaluated states: %v | States per second: %v\n", trackingData.InitializedBoardStates, trackingData.EvaluatedBoardStates, 1000000000*trackingData.EvaluatedBoardStates/uint(trackingData.Duration))
	fmt.Println("\n--------PATH--------\n")

	var previousBoardState types.BoardState
	var previousRobotOrder = robotOrder
	for indexBoardState, boardState := range path {
		robots := helper.SeparateRobots(boardState)
		cRobotOrder, err := getNewRobotOrder(previousBoardState, previousRobotOrder, boardState)
		if err != nil {
			return err
		}

		if indexBoardState == 0 {
			fmt.Printf("Start\n")

			previousBoardState = boardState

			for indexRobot, robotPosition := range robots {
				robotColor, err := getRobotColorByIndex(robotOrder, uint8(indexRobot))
				if err != nil {
					return err
				}
				msg := fmt.Sprintf("%+v ", helper.ConvBytePositionToPosition(byte(robotPosition)))

				switch robotColor {
				case "yellow":
					color.HiYellow(msg)
				case "red":
					color.HiRed(msg)
				case "green":
					color.HiGreen(msg)
				case "blue":
					color.HiBlue(msg)
				}
			}

			previousRobotOrder = robotOrder
		} else {
			robotColor, direction, err := getMovedRobotColorAndDirection(previousBoardState, boardState, cRobotOrder)
			if err != nil {
				return err
			}

			fmt.Printf("Move: %v", indexBoardState)
			if indexBoardState < 10 {
				fmt.Printf("\t | ")
			} else {
				fmt.Printf(" | ")
			}

			switch robotColor {
			case "yellow":
				color.HiYellow(direction)
			case "red":
				color.HiRed(direction)
			case "green":
				color.HiGreen(direction)
			case "blue":
				color.HiBlue(direction)
			}

			previousBoardState = boardState
			previousRobotOrder = cRobotOrder
		}

		if indexBoardState == (len(path) - 1) {
			fmt.Printf("\nFinish\n")

			for indexRobot, robotPosition := range robots {
				robotColor, err := getRobotColorByIndex(cRobotOrder, uint8(indexRobot))
				if err != nil {
					return err
				}
				msg := fmt.Sprintf("%+v ", helper.ConvBytePositionToPosition(byte(robotPosition)))

				switch robotColor {
				case "yellow":
					color.HiYellow(msg)
				case "red":
					color.HiRed(msg)
				case "green":
					color.HiGreen(msg)
				case "blue":
					color.HiBlue(msg)
				}
			}
		}

		fmt.Println("")
	}
	fmt.Printf("====================\n")
	return nil
}

// RobotStoppingPositions Prints per node the stopping positions of a theoretical moved robot per direction
func RobotStoppingPositions(robotStoppingPositions *types.RobotStoppingPositions) (err error) {
	fmt.Printf("\n\n====================\n")
	for nodePosition, stoppingPositions := range *robotStoppingPositions {
		fmt.Printf("Node: %v\t| ", nodePosition)
		for direction, stoppingPosition := range stoppingPositions {
			fmt.Printf("%v: %v\t", direction, stoppingPosition)
		}
		fmt.Println("")
	}
	fmt.Printf("====================\n")
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

// getRobotColorByIndex Returns the color name of the robot at the given index
func getRobotColorByIndex(robotOrder types.RobotOrder, index uint8) (string, error) {
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

	robotColor, err = getRobotColorByIndex(currentRobotOrder, robotIndex)
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
