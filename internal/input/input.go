package input

import (
	"Ricochet-Robot-Solver/internal/helper"
	"Ricochet-Robot-Solver/internal/precomputation"
	"Ricochet-Robot-Solver/internal/types"
	"encoding/json"
	"os"
	"sort"
)

// GetData Manages the complete data loading process and returns the formatted data
func GetData(boardDataLocation string) (board types.Board, initBoardState types.BoardState, initRobotOrder types.RobotOrder, robotStoppingPositions types.RobotStoppingPositions, err error) {
	data, err := getJsonData(boardDataLocation)
	if err != nil {
		return types.Board{}, 0, 0, types.RobotStoppingPositions{}, err
	}

	board, initRobotPositions, err := loadData(data)
	if err != nil {
		return types.Board{}, 0, 0, types.RobotStoppingPositions{}, err
	}

	initBoardState, initRobotOrder, err = getInitBoardState(&board, initRobotPositions)
	if err != nil {
		return types.Board{}, 0, 0, types.RobotStoppingPositions{}, err
	}

	robotStoppingPositions, err = precomputation.PrecomputeRobotMoves(&board)
	if err != nil {
		return types.Board{}, 0, 0, types.RobotStoppingPositions{}, err
	}

	return board, initBoardState, initRobotOrder, robotStoppingPositions, nil
}

// getInitBoardState Builds the initial board state based on the board configuration
func getInitBoardState(board *types.Board, initRobotPositions [4]byte) (initBoardState types.BoardState, initRobotOrder types.RobotOrder, err error) {
	targetColor, err := helper.GetTargetColor(board.Target)
	if err != nil {
		return 0, 0, err
	}

	robots := initRobotPositions

	var robotsSlice = robots[0:4]
	sort.Slice(robotsSlice, func(i, j int) bool {
		return robotsSlice[i] < robotsSlice[j]
	})

	var colors = []string{"yellow", "red", "green", "blue"}

	switch targetColor {
	case "yellow":
		colors = append(colors[:0], colors[1:]...)
	case "red":
		colors = append(colors[:1], colors[2:]...)
	case "green":
		colors = append(colors[:2], colors[3:]...)
	case "blue":
		colors = append(colors[:3], colors[4:]...)
	}

	helper.SetRobotColorByIndex(&initRobotOrder, targetColor, 0)

	for _, color := range colors {
		for indexRobotSlice, robotSlice := range robotsSlice {
			switch color {
			case "yellow":
				if robotSlice == initRobotPositions[0] {
					helper.SetRobotColorByIndex(&initRobotOrder, "yellow", uint8(indexRobotSlice))
				}
			case "red":
				if robotSlice == initRobotPositions[1] {
					helper.SetRobotColorByIndex(&initRobotOrder, "red", uint8(indexRobotSlice))
				}
			case "green":
				if robotSlice == initRobotPositions[2] {
					helper.SetRobotColorByIndex(&initRobotOrder, "green", uint8(indexRobotSlice))
				}
			case "blue":
				if robotSlice == initRobotPositions[3] {
					helper.SetRobotColorByIndex(&initRobotOrder, "blue", uint8(indexRobotSlice))
				}
			}
		}
	}

	initBoardState = types.BoardState(uint32(initRobotPositions[helper.GetRobotColorCodeByIndex(initRobotOrder, 0)])<<24 | uint32(initRobotPositions[helper.GetRobotColorCodeByIndex(initRobotOrder, 1)])<<16 | uint32(initRobotPositions[helper.GetRobotColorCodeByIndex(initRobotOrder, 2)])<<8 | uint32(initRobotPositions[helper.GetRobotColorCodeByIndex(initRobotOrder, 3)])<<0)
	return initBoardState, initRobotOrder, nil
}

// getJsonDate Loads the json date into a byte array
func getJsonData(path string) (jsonData []byte, err error) {
	jsonData, err = os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

// loadData Loads the json byte array into the predefined data structure
func loadData(data []byte) (board types.Board, initRobotPositions [4]byte, err error) {
	var rawBoard types.RawBoard
	err = json.Unmarshal(data, &rawBoard)

	if err != nil {
		return types.Board{}, [4]byte{}, err
	}

	// convert board_data
	board, initRobotPositions, err = convData(rawBoard)
	if err != nil {
		return types.Board{}, [4]byte{}, err
	}

	// precompute the minimal moves that a robot has to make in order to get to the target node, if it could stop wherever it wants
	err = precomputation.PrecomputeMinimalMoveCounts(&board)

	return board, initRobotPositions, nil
}

// convData Translates the raw data into the optimized data structure
func convData(data types.RawBoard) (board types.Board, initRobotPositions [4]byte, err error) {
	// node conversion
	for rowIndex := 0; rowIndex < 16; rowIndex++ {
		for columnIndex := 0; columnIndex < 16; columnIndex++ {
			var node types.Node

			// distance between node and target
			helper.SetBit((*byte)(&node), 7, true)
			helper.SetBit((*byte)(&node), 6, true)
			helper.SetBit((*byte)(&node), 5, true)

			helper.SetBit((*byte)(&node), 4, false) // is a robot present

			// neighbors
			helper.SetBit((*byte)(&node), 3, true)
			helper.SetBit((*byte)(&node), 2, true)
			helper.SetBit((*byte)(&node), 1, true)
			helper.SetBit((*byte)(&node), 0, true)

			board.Grid[rowIndex][columnIndex] = node
		}
	}

	// set walls at edges
	for _, edgeDirection := range [4]string{"top", "bottom", "left", "right"} {
		switch edgeDirection {
		case "top":
			rowIndex := 0
			for columnIndex := 0; columnIndex < 16; columnIndex++ {
				helper.SetBit((*byte)(&(board.Grid[rowIndex][columnIndex])), 3, false)
			}
		case "bottom":
			rowIndex := 15
			for columnIndex := 0; columnIndex < 16; columnIndex++ {
				helper.SetBit((*byte)(&(board.Grid[rowIndex][columnIndex])), 2, false)
			}
		case "left":
			columnIndex := 0
			for rowIndex := 0; rowIndex < 16; rowIndex++ {
				helper.SetBit((*byte)(&(board.Grid[rowIndex][columnIndex])), 1, false)
			}
		case "right":
			columnIndex := 15
			for rowIndex := 0; rowIndex < 16; rowIndex++ {
				helper.SetBit((*byte)(&(board.Grid[rowIndex][columnIndex])), 0, false)
			}
		}
	}

	// add walls to the board -> set the designated  bits of the nodes
	for _, wall := range data.Walls {

		switch wall.Direction1 {
		case "top":
			helper.SetBit((*byte)(&(board.Grid[wall.Position1.Row][wall.Position1.Column])), 3, false)
		case "bottom":
			helper.SetBit((*byte)(&(board.Grid[wall.Position1.Row][wall.Position1.Column])), 2, false)
		case "left":
			helper.SetBit((*byte)(&(board.Grid[wall.Position1.Row][wall.Position1.Column])), 1, false)
		case "right":
			helper.SetBit((*byte)(&(board.Grid[wall.Position1.Row][wall.Position1.Column])), 0, false)
		}

		switch wall.Direction2 {
		case "top":
			helper.SetBit((*byte)(&(board.Grid[wall.Position2.Row][wall.Position2.Column])), 3, false)
		case "bottom":
			helper.SetBit((*byte)(&(board.Grid[wall.Position2.Row][wall.Position2.Column])), 2, false)
		case "left":
			helper.SetBit((*byte)(&(board.Grid[wall.Position2.Row][wall.Position2.Column])), 1, false)
		case "right":
			helper.SetBit((*byte)(&(board.Grid[wall.Position2.Row][wall.Position2.Column])), 0, false)
		}
	}

	// target conversion
	var targetColorAndSymbol byte
	targetPosition := helper.ConvPositionToByte(data.Target.Position)

	switch data.Target.Color {
	case "yellow":
		helper.SetBit(&targetColorAndSymbol, 7, true)
	case "red":
		helper.SetBit(&targetColorAndSymbol, 6, true)
	case "green":
		helper.SetBit(&targetColorAndSymbol, 5, true)
	case "blue":
		helper.SetBit(&targetColorAndSymbol, 4, true)
	}

	switch data.Target.Symbol {
	case "circle":
		helper.SetBit(&targetColorAndSymbol, 3, true)
	case "triangle":
		helper.SetBit(&targetColorAndSymbol, 2, true)
	case "square":
		helper.SetBit(&targetColorAndSymbol, 1, true)
	case "hexagon":
		helper.SetBit(&targetColorAndSymbol, 0, true)
	}

	board.Target = types.Target(uint16(targetColorAndSymbol)<<8 | uint16(targetPosition))

	// Robot conversion
	for colorIndex, color := range [4]string{"yellow", "red", "green", "blue"} {
		for _, robot := range data.Robots {
			if color == robot.Color {
				initRobotPositions[colorIndex] = helper.ConvPositionToByte(robot.Position)
			}
		}
	}

	return board, initRobotPositions, nil
}
