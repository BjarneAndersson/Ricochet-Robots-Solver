package output

import (
	"../helper"
	"../tracker"
	"../types"
	"fmt"
	"github.com/fatih/color"
	"strconv"
)

func BoardState(boardState types.BoardState, robotColors types.RobotColors) (err error) {
	fmt.Printf("\n\n====================\n")
	fmt.Printf("Board State: %v | %v\n", boardState, convertNumberToBits(int(boardState)))

	robots := helper.SeparateRobots(boardState)
	for indexRobot, robotPosition := range robots {
		robotColor, err := helper.GetRobotColorByIndex(robotColors, uint8(indexRobot))
		if err != nil {
			return err
		}
		msg := fmt.Sprintf("%+v\n", helper.ConvBytePositionToPosition(robotPosition))

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
	fmt.Printf("====================\n")
	return nil
}

func Path(path []types.BoardState, trackingData tracker.TrackingDataSolver, robotColors types.RobotColors) (err error) {
	fmt.Printf("\n\n====================\n")
	fmt.Printf("Moves\t%d\n", len(path)-1)
	fmt.Printf("Time\t%s\n", trackingData.Duration)
	fmt.Printf("Path\t%+v\n", path)
	fmt.Printf("Initialized states: %v | Evaluated states: %v\n", trackingData.InitializedBoardStates, trackingData.EvaluatedBoardStates)
	fmt.Println("")

	var previousBoardState types.BoardState
	for indexBoardState, boardState := range path {
		robots := helper.SeparateRobots(boardState)

		switch indexBoardState {
		case 0:
			fmt.Printf("Start\t| ")
			for indexRobot, robotPosition := range robots {
				robotColor, err := helper.GetRobotColorByIndex(robotColors, uint8(indexRobot))
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
			robotColor, direction, err := getMovedRobotColorAndDirection(previousBoardState, boardState, robotColors)
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
				robotColor, err := helper.GetRobotColorByIndex(robotColors, uint8(indexRobot))
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
			robotColor, direction, err := getMovedRobotColorAndDirection(previousBoardState, boardState, robotColors)
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
				robotColor, err := helper.GetRobotColorByIndex(robotColors, uint8(indexRobot))
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

func convertNumberToBits(number int) string {
	return strconv.FormatInt(int64(number), 2)
}

func getMovedRobotColorAndDirection(previousBoardState types.BoardState, currentBoardState types.BoardState, robotColors types.RobotColors) (robotColor string, direction string, err error) {
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

	robotColor, err = helper.GetRobotColorByIndex(robotColors, robotIndex)
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
