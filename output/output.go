package output

import (
	"../helper"
	"../types"
	"fmt"
	"strconv"
)

func BoardState(boardState types.BoardState, robotColors types.RobotColors) (err error) {
	fmt.Printf("\n\n====================\n")
	fmt.Printf("Board State: %v | %v\n", boardState, convertNumberToBits(int(boardState)))

	fmt.Printf("Robots:\n")
	robots := helper.SeparateRobots(boardState)
	for indexRobot, robotPosition := range robots {
		robotColor, err := helper.GetRobotColorByIndex(robotColors, uint8(indexRobot))
		if err != nil {
			return err
		}
		fmt.Printf("%v: %+v\n", robotColor, helper.ConvBytePositionToPosition(robotPosition))
	}
	fmt.Printf("====================\n")
	return nil
}

func convertNumberToBits(number int) string {
	return strconv.FormatInt(int64(number), 2)
}
