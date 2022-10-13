package output

import (
	"../types"
	"sort"
	"testing"
)

func TestCalcNewRobotOrder(t *testing.T) {
	previousRobotOrder := createRobotOrder([4]types.RobotColor{3, 2, 1, 0})
	previousBoardState := createNewBoardState([4]byte{175, 22, 53, 117})
	nextBoardState := createNewBoardState([4]byte{175, 11, 22, 53})

	nextRobotOrder := createRobotOrder([4]types.RobotColor{3, 0, 2, 1})

	result, err := GetNewRobotOrder(previousBoardState, previousRobotOrder, nextBoardState)
	if err != nil {
		t.Errorf("CreateNewRobotOrder() FAILED. Error: %v", err)
	}

	if result != nextRobotOrder {
		t.Errorf("CreateNewRobotOrder() FAILED. Expected %v, got %v", nextRobotOrder, result)
	} else {
		t.Logf("CreateNewRobotOrder() PASSED. Expected %v, got %v", nextRobotOrder, result)
	}
}

func createRobotOrder(colors [4]types.RobotColor) byte {
	var robotOrder byte = 0
	for index, color := range colors {
		switch index {
		case 0:
			robotOrder = robotOrder | (byte(color) << 6)
		case 1:
			robotOrder = robotOrder | (byte(color) << 4)
		case 2:
			robotOrder = robotOrder | (byte(color) << 2)
		case 3:
			robotOrder = robotOrder | (byte(color) << 0)
		}
	}
	return robotOrder
}

func createNewBoardState(robots [4]byte) types.BoardState {
	tRobots := robots

	var robotsSlice = tRobots[1:4]
	sort.Slice(robotsSlice, func(i, j int) bool {
		return robotsSlice[i] < robotsSlice[j]
	})
	return types.BoardState(uint32(robots[0])<<24 | uint32(robotsSlice[0])<<16 | uint32(robotsSlice[1])<<8 | uint32(robotsSlice[2])<<0)
}
