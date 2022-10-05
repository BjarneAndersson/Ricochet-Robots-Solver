package input

import (
	"../bitOperations"
	"../helper"
	"../precomputation"
	"../types"
	"encoding/json"
	"os"
)

func GetData(boardDataLocation string) (gameRound types.GameRound, initBoardState types.BoardState, robotStoppingPositions types.RobotStoppingPositions, err error) {
	data, err := getJsonData(boardDataLocation)
	if err != nil {
		return types.GameRound{}, 0, types.RobotStoppingPositions{}, err
	}

	gameRound, err = loadData(data)
	if err != nil {
		return types.GameRound{}, 0, types.RobotStoppingPositions{}, err
	}

	initBoardState, err = getInitBoardState(&gameRound)
	if err != nil {
		return types.GameRound{}, 0, types.RobotStoppingPositions{}, err
	}

	robotStoppingPositions, err = precomputation.PrecomputeRobotMoves(&gameRound)
	if err != nil {
		return types.GameRound{}, 0, types.RobotStoppingPositions{}, err
	}

	return gameRound, initBoardState, robotStoppingPositions, nil
}

func getInitBoardState(gameRound *types.GameRound) (initBoardState types.BoardState, err error) {

	targetColor, err := helper.GetTargetColor(gameRound.Target)
	var targetColorBaseIndex uint8
	if err != nil {
		return 0, err
	}

	switch targetColor {
	case "yellow":
		targetColorBaseIndex = 0
		gameRound.RobotColors["yellow"] = 0
	case "red":
		targetColorBaseIndex = 1
		gameRound.RobotColors["yellow"] = 1
		gameRound.RobotColors["red"] = 0
	case "green":
		targetColorBaseIndex = 2
		gameRound.RobotColors["yellow"] = 1
		gameRound.RobotColors["red"] = 2
		gameRound.RobotColors["green"] = 0
	case "blue":
		targetColorBaseIndex = 3
		gameRound.RobotColors["yellow"] = 1
		gameRound.RobotColors["red"] = 2
		gameRound.RobotColors["green"] = 3
		gameRound.RobotColors["blue"] = 0
	}

	robotsToSort := []uint8{0, 1, 2, 3}

	robotsToSort = append(robotsToSort[:targetColorBaseIndex], robotsToSort[targetColorBaseIndex+1:]...)

	initBoardState = types.BoardState(uint32(gameRound.Robots[targetColorBaseIndex])<<24 | uint32(gameRound.Robots[robotsToSort[0]])<<16 | uint32(gameRound.Robots[robotsToSort[1]])<<8 | uint32(gameRound.Robots[robotsToSort[2]])<<0)

	return initBoardState, nil
}

func getJsonData(path string) (jsonData []byte, err error) {
	jsonData, err = os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func loadData(data []byte) (gameRound types.GameRound, err error) {
	var rawBoard types.RawBoard
	err = json.Unmarshal(data, &rawBoard)

	if err != nil {
		return types.GameRound{}, err
	}

	// convert board_data
	gameRound, err = convData(rawBoard)
	if err != nil {
		return types.GameRound{}, err
	}

	err = precomputation.PrecomputeBoard(&gameRound)

	return gameRound, nil
}

// convData convert board_data in json format to board_data byte format
func convData(data types.RawBoard) (gameRound types.GameRound, err error) {
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

			gameRound.Board[rowIndex][columnIndex] = cell
		}
	}

	// set walls at edges
	for _, edgeDirection := range [4]string{"top", "bottom", "left", "right"} {
		switch edgeDirection {
		case "top":
			rowIndex := 0
			for columnIndex := 0; columnIndex < 16; columnIndex++ {
				bitOperations.SetBit(&(gameRound.Board[rowIndex][columnIndex]), 3, false)
			}
		case "bottom":
			rowIndex := 15
			for columnIndex := 0; columnIndex < 16; columnIndex++ {
				bitOperations.SetBit(&(gameRound.Board[rowIndex][columnIndex]), 2, false)
			}
		case "left":
			columnIndex := 0
			for rowIndex := 0; rowIndex < 16; rowIndex++ {
				bitOperations.SetBit(&(gameRound.Board[rowIndex][columnIndex]), 1, false)
			}
		case "right":
			columnIndex := 15
			for rowIndex := 0; rowIndex < 16; rowIndex++ {
				bitOperations.SetBit(&(gameRound.Board[rowIndex][columnIndex]), 0, false)
			}
		}
	}

	// add walls to the gameRound
	for _, wall := range data.Walls {

		switch wall.Direction1 {
		case "top":
			bitOperations.SetBit(&(gameRound.Board[wall.Position1.Row][wall.Position1.Column]), 3, false)
		case "bottom":
			bitOperations.SetBit(&(gameRound.Board[wall.Position1.Row][wall.Position1.Column]), 2, false)
		case "left":
			bitOperations.SetBit(&(gameRound.Board[wall.Position1.Row][wall.Position1.Column]), 1, false)
		case "right":
			bitOperations.SetBit(&(gameRound.Board[wall.Position1.Row][wall.Position1.Column]), 0, false)
		}

		switch wall.Direction2 {
		case "top":
			bitOperations.SetBit(&(gameRound.Board[wall.Position2.Row][wall.Position2.Column]), 3, false)
		case "bottom":
			bitOperations.SetBit(&(gameRound.Board[wall.Position2.Row][wall.Position2.Column]), 2, false)
		case "left":
			bitOperations.SetBit(&(gameRound.Board[wall.Position2.Row][wall.Position2.Column]), 1, false)
		case "right":
			bitOperations.SetBit(&(gameRound.Board[wall.Position2.Row][wall.Position2.Column]), 0, false)
		}
	}

	// target conversion
	var targetColorAndSymbol byte
	targetPosition := helper.ConvPosToByte(data.Target.Position)

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

	gameRound.Target = uint16(targetColorAndSymbol)<<8 | uint16(targetPosition)

	// Robot conversion
	gameRound.RobotColors = types.RobotColors{
		"yellow": 0,
		"red":    1,
		"green":  2,
		"blue":   3,
	}

	for indexRobot, robot := range data.Robots {
		gameRound.Robots[gameRound.RobotColors[robot.Color]] = helper.ConvPosToByte(data.Robots[indexRobot].Position)
	}

	return gameRound, nil
}
