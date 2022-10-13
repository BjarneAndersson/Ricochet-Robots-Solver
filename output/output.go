package output

import (
	"../helper"
	"../tracker"
	"../types"
	"fmt"
	"github.com/fatih/color"
)

func Neighbors(gameRound *types.GameRound) (err error) {
	fmt.Printf("\n\n====================\n")
	for rowIndex := range gameRound.Board {
		for _, node := range gameRound.Board[rowIndex] {
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

func BoardState(boardState types.BoardState, trackingData tracker.TrackingDataSolver) (err error) {
	boardStateInBits, err := convertNumberToBits(int(boardState), 32)
	if err != nil {
		return err
	}

	fmt.Printf("\n\n====================\n")
	fmt.Printf("GameRound State: %v | %v\n", boardState, boardStateInBits)

	robots := helper.SeparateRobots(boardState)
	for _, robotPosition := range robots {
		fmt.Printf("%+v\n", helper.ConvBytePositionToPosition(robotPosition))
	}
	fmt.Printf("Initialized states: %v | Evaluated states: %v\n", trackingData.InitializedBoardStates, trackingData.EvaluatedBoardStates)
	fmt.Printf("====================\n")
	return nil
}

func Path(path []types.BoardState, trackingData tracker.TrackingDataSolver, robotOrder byte) (err error) {
	fmt.Printf("\n\n====================\n")
	fmt.Printf("Moves\t%d\n", len(path)-1)
	fmt.Printf("Time\t%s\n", trackingData.Duration)
	fmt.Printf("Path\t%+v\n", path)
	fmt.Printf("Initialized states: %v | Evaluated states: %v | States per second: %v\n", trackingData.InitializedBoardStates, trackingData.EvaluatedBoardStates, 100000000*trackingData.EvaluatedBoardStates/uint(trackingData.Duration))
	fmt.Println("")

	var previousBoardState types.BoardState
	var previousRobotOrder byte = robotOrder
	for indexBoardState, boardState := range path {
		robots := helper.SeparateRobots(boardState)
		cRobotOrder, err := getNewRobotOrder(previousBoardState, previousRobotOrder, boardState)
		if err != nil {
			return err
		}

		switch indexBoardState {
		case 0:
			fmt.Printf("Start\t| ")
			for indexRobot, robotPosition := range robots {
				robotColor, err := getRobotColorByIndex(robotOrder, uint8(indexRobot))
				if err != nil {
					return err
				}
				msg := fmt.Sprintf("%+v ", helper.ConvBytePositionToPosition(robotPosition))

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

			previousBoardState = boardState
		case len(path) - 1:
			robotColor, direction, err := getMovedRobotColorAndDirection(previousBoardState, boardState, cRobotOrder)
			if err != nil {
				return err
			}

			fmt.Println("")

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

			fmt.Println("")

			previousBoardState = boardState

			fmt.Printf("\nFinish\t| ")

			for indexRobot, robotPosition := range robots {
				robotColor, err := getRobotColorByIndex(cRobotOrder, uint8(indexRobot))
				if err != nil {
					return err
				}
				msg := fmt.Sprintf("%+v ", helper.ConvBytePositionToPosition(robotPosition))

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
		default:
			robotColor, direction, err := getMovedRobotColorAndDirection(previousBoardState, boardState, cRobotOrder)
			if err != nil {
				return err
			}

			fmt.Println("")

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

			fmt.Println("")

			previousBoardState = boardState

			fmt.Printf("\nMove: %v\t| ", indexBoardState)

			for indexRobot, robotPosition := range robots {
				robotColor, err := getRobotColorByIndex(cRobotOrder, uint8(indexRobot))
				if err != nil {
					return err
				}
				msg := fmt.Sprintf("%+v ", helper.ConvBytePositionToPosition(robotPosition))

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

		if indexBoardState != 0 {

		} else if indexBoardState == 0 {
			previousBoardState = boardState
		}
		fmt.Println("")
	}
	fmt.Printf("====================\n")
	return nil
}

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

func getRobotColorByIndex(robotOrder byte, index uint8) (string, error) {
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

func getMovedRobotColorAndDirection(previousBoardState types.BoardState, currentBoardState types.BoardState, robotOrder byte) (robotColor string, direction string, err error) {
	preRobots := helper.SeparateRobots(previousBoardState)
	curRobots := helper.SeparateRobots(currentBoardState)

	var robotIndex uint8

	for indexPreRobot, preRobot := range preRobots {
		if preRobot == curRobots[indexPreRobot] {
			continue
		} else {
			robotIndex = uint8(indexPreRobot)
		}
	}

	robotColor, err = getRobotColorByIndex(robotOrder, robotIndex)
	if err != nil {
		return "", "", err
	}

	direction, err = evaluateDirectionChange(preRobots[robotIndex], curRobots[robotIndex])
	if err != nil {
		return "", "", err
	}

	return robotColor, direction, nil
}

func evaluateDirectionChange(previousRobot byte, currentRobot byte) (direction string, err error) {
	preRobotPos := helper.ConvBytePositionToPosition(previousRobot)
	curRobotPos := helper.ConvBytePositionToPosition(currentRobot)

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

func getNewRobotOrder(previousBoardState types.BoardState, previousRobotOrder byte, currentBoardState types.BoardState) (currentRobotOrder byte, err error) {
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

func robotInRobots(robot byte, robots [4]byte) (exist bool, index int) {
	for iterIndexRobot, iterRobot := range robots {
		if iterRobot == robot {
			return true, iterIndexRobot
		}
	}
	return false, -1
}
