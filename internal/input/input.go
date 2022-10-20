package input

import (
	"Ricochet-Robot-Solver/internal/bitOperations"
	"Ricochet-Robot-Solver/internal/helper"
	"Ricochet-Robot-Solver/internal/precomputation"
	"Ricochet-Robot-Solver/internal/types"
	"encoding/json"
	"os"
)

func GetData(boardDataLocation string) (board types.Board, initBoardState types.BoardState, robotStoppingPositions types.RobotStoppingPositions, err error) {
	data, err := getJsonData(boardDataLocation)
	if err != nil {
		return types.Board{}, 0, types.RobotStoppingPositions{}, err
	}

	board, err = loadData(data)
	if err != nil {
		return types.Board{}, 0, types.RobotStoppingPositions{}, err
	}

	initBoardState, err = getInitBoardState(&board)
	if err != nil {
		return types.Board{}, 0, types.RobotStoppingPositions{}, err
	}

	robotStoppingPositions, err = precomputation.PrecomputeRobotMoves(&board)
	if err != nil {
		return types.Board{}, 0, types.RobotStoppingPositions{}, err
	}

	return board, initBoardState, robotStoppingPositions, nil
}

func getInitBoardState(board *types.Board) (initBoardState types.BoardState, err error) {

	targetColor, err := helper.GetTargetColor(board.Target)
	var targetColorBaseIndex uint8
	if err != nil {
		return 0, err
	}

	switch targetColor {
	case "yellow":
		targetColorBaseIndex = 0
		board.RobotColors["yellow"] = 0
	case "red":
		targetColorBaseIndex = 1
		board.RobotColors["yellow"] = 1
		board.RobotColors["red"] = 0
	case "green":
		targetColorBaseIndex = 2
		board.RobotColors["yellow"] = 1
		board.RobotColors["red"] = 2
		board.RobotColors["green"] = 0
	case "blue":
		targetColorBaseIndex = 3
		board.RobotColors["yellow"] = 1
		board.RobotColors["red"] = 2
		board.RobotColors["green"] = 3
		board.RobotColors["blue"] = 0
	}

	robotsToSort := []uint8{0, 1, 2, 3}

	robotsToSort = append(robotsToSort[:targetColorBaseIndex], robotsToSort[targetColorBaseIndex+1:]...)

	initBoardState = types.BoardState(uint32(board.Robots[targetColorBaseIndex])<<24 | uint32(board.Robots[robotsToSort[0]])<<16 | uint32(board.Robots[robotsToSort[1]])<<8 | uint32(board.Robots[robotsToSort[2]])<<0)

	return initBoardState, nil
}

func getJsonData(path string) (jsonData []byte, err error) {
	jsonData, err = os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func loadData(data []byte) (board types.Board, err error) {
	var rawBoard types.RawBoard
	err = json.Unmarshal(data, &rawBoard)

	if err != nil {
		return types.Board{}, err
	}

	// convert board_data
	board, err = convData(rawBoard)
	if err != nil {
		return types.Board{}, err
	}

	err = precomputation.PrecomputeBoard(&board)

	return board, nil
}

// convData convert board_data in json format to board_data byte format
func convData(data types.RawBoard) (board types.Board, err error) {
	// node conversion
	for rowIndex := 0; rowIndex < 16; rowIndex++ {
		for columnIndex := 0; columnIndex < 16; columnIndex++ {
			var cell byte

			// distance between node and target
			bitOperations.SetBit(&cell, 7, true)
			bitOperations.SetBit(&cell, 6, true)
			bitOperations.SetBit(&cell, 5, true)

			bitOperations.SetBit(&cell, 4, false) // is a robot present

			// neighbors
			bitOperations.SetBit(&cell, 3, true)
			bitOperations.SetBit(&cell, 2, true)
			bitOperations.SetBit(&cell, 1, true)
			bitOperations.SetBit(&cell, 0, true)

			board.Board[rowIndex][columnIndex] = cell
		}
	}

	// set walls at edges
	for _, edgeDirection := range [4]string{"top", "bottom", "left", "right"} {
		switch edgeDirection {
		case "top":
			rowIndex := 0
			for columnIndex := 0; columnIndex < 16; columnIndex++ {
				bitOperations.SetBit(&(board.Board[rowIndex][columnIndex]), 3, false)
			}
		case "bottom":
			rowIndex := 15
			for columnIndex := 0; columnIndex < 16; columnIndex++ {
				bitOperations.SetBit(&(board.Board[rowIndex][columnIndex]), 2, false)
			}
		case "left":
			columnIndex := 0
			for rowIndex := 0; rowIndex < 16; rowIndex++ {
				bitOperations.SetBit(&(board.Board[rowIndex][columnIndex]), 1, false)
			}
		case "right":
			columnIndex := 15
			for rowIndex := 0; rowIndex < 16; rowIndex++ {
				bitOperations.SetBit(&(board.Board[rowIndex][columnIndex]), 0, false)
			}
		}
	}

	// add walls to the board
	for _, wall := range data.Walls {

		switch wall.Direction1 {
		case "top":
			bitOperations.SetBit(&(board.Board[wall.Position1.Row][wall.Position1.Column]), 3, false)
		case "bottom":
			bitOperations.SetBit(&(board.Board[wall.Position1.Row][wall.Position1.Column]), 2, false)
		case "left":
			bitOperations.SetBit(&(board.Board[wall.Position1.Row][wall.Position1.Column]), 1, false)
		case "right":
			bitOperations.SetBit(&(board.Board[wall.Position1.Row][wall.Position1.Column]), 0, false)
		}

		switch wall.Direction2 {
		case "top":
			bitOperations.SetBit(&(board.Board[wall.Position2.Row][wall.Position2.Column]), 3, false)
		case "bottom":
			bitOperations.SetBit(&(board.Board[wall.Position2.Row][wall.Position2.Column]), 2, false)
		case "left":
			bitOperations.SetBit(&(board.Board[wall.Position2.Row][wall.Position2.Column]), 1, false)
		case "right":
			bitOperations.SetBit(&(board.Board[wall.Position2.Row][wall.Position2.Column]), 0, false)
		}
	}

	// target conversion
	var targetColorAndSymbol byte
	var targetPosition byte

	switch data.Target.Color {
	case "yellow":
		bitOperations.SetBit(&targetColorAndSymbol, 7, true)
	case "red":
		bitOperations.SetBit(&targetColorAndSymbol, 6, true)
	case "green":
		bitOperations.SetBit(&targetColorAndSymbol, 5, true)
	case "blue":
		bitOperations.SetBit(&targetColorAndSymbol, 4, true)
	}

	switch data.Target.Symbol {
	case "circle":
		bitOperations.SetBit(&targetColorAndSymbol, 3, true)
	case "triangle":
		bitOperations.SetBit(&targetColorAndSymbol, 2, true)
	case "square":
		bitOperations.SetBit(&targetColorAndSymbol, 1, true)
	case "hexagon":
		bitOperations.SetBit(&targetColorAndSymbol, 0, true)
	}

	helper.ConvPosToByte(&targetPosition, data.Target.Position.Column, data.Target.Position.Row)

	board.Target = uint16(targetColorAndSymbol)<<8 | uint16(targetPosition)

	// Robot conversion
	board.RobotColors = types.RobotColors{
		"yellow": 0,
		"red":    1,
		"green":  2,
		"blue":   3,
	}

	for indexRobot, robot := range data.Robots {
		switch robot.Color {
		case "yellow":
			helper.ConvPosToByte(&board.Robots[board.RobotColors[robot.Color]], data.Robots[indexRobot].Position.Column, data.Robots[indexRobot].Position.Row)
		case "red":
			helper.ConvPosToByte(&board.Robots[board.RobotColors[robot.Color]], data.Robots[indexRobot].Position.Column, data.Robots[indexRobot].Position.Row)
		case "green":
			helper.ConvPosToByte(&board.Robots[board.RobotColors[robot.Color]], data.Robots[indexRobot].Position.Column, data.Robots[indexRobot].Position.Row)
		case "blue":
			helper.ConvPosToByte(&board.Robots[board.RobotColors[robot.Color]], data.Robots[indexRobot].Position.Column, data.Robots[indexRobot].Position.Row)
		}
	}

	return board, nil
}
