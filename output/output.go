package output

import (
	"../helper"
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

func convertNumberToBits(number int) string {
	return strconv.FormatInt(int64(number), 2)
}
