package input

import (
	"../bitOperations"
	"../helper"
	"../precomputation"
	"../types"
	"encoding/json"
	"os"
)

func GetData() (board types.Board, initBoardState uint32, err error) {
	data, err := getJsonData("K:\\Coding\\Golang\\ricochet-robot-solver\\board_data.json")
	if err != nil {
		return types.Board{}, 0, err
	}

	board, err = loadData(data)
	if err != nil {
		return types.Board{}, 0, err
	}

	initBoardState, err = getInitBoardState(&board)
	if err != nil {
		return types.Board{}, 0, err
	}

	return board, initBoardState, nil
}

func getInitBoardState(board *types.Board) (initBoardState uint32, err error) {

	targetColor, err := helper.GetTargetColor(board.Target)
	if err != nil {
		return 0, err
	}

	robotsToSort := []uint8{0, 1, 2, 3}

	robotsToSort = append(robotsToSort[:targetColor], robotsToSort[targetColor+1:]...)

	initBoardState = uint32(board.Robots[targetColor])<<24 | uint32(board.Robots[robotsToSort[0]])<<16 | uint32(board.Robots[robotsToSort[1]])<<8 | uint32(board.Robots[robotsToSort[2]])<<0

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
	for rowIndex, rowOfNodes := range data.Nodes {
		for columnIndex, node := range rowOfNodes {
			var cell byte

			// distance between node and target
			bitOperations.SetBit(&cell, 7, true)
			bitOperations.SetBit(&cell, 6, true)
			bitOperations.SetBit(&cell, 5, true)

			bitOperations.SetBit(&cell, 4, false) // is a robot present

			bitOperations.SetBit(&cell, 3, node.Neighbors.Up)
			bitOperations.SetBit(&cell, 2, node.Neighbors.Down)
			bitOperations.SetBit(&cell, 1, node.Neighbors.Left)
			bitOperations.SetBit(&cell, 0, node.Neighbors.Right)

			board.Board[rowIndex][columnIndex] = cell
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

	convPosToByte(&targetPosition, data.Target.Position.Column, data.Target.Position.Row)

	board.Target = uint16(targetColorAndSymbol)<<8 | uint16(targetPosition)

	// Robot conversion

	for indexRobot, robot := range data.Robots {
		switch robot.Color {
		case "yellow":
			convPosToByte(&board.Robots[types.Yellow], data.Robots[indexRobot].Position.Column, data.Robots[indexRobot].Position.Row)
		case "red":
			convPosToByte(&board.Robots[types.Red], data.Robots[indexRobot].Position.Column, data.Robots[indexRobot].Position.Row)
		case "green":
			convPosToByte(&board.Robots[types.Green], data.Robots[indexRobot].Position.Column, data.Robots[indexRobot].Position.Row)
		case "blue":
			convPosToByte(&board.Robots[types.Blue], data.Robots[indexRobot].Position.Column, data.Robots[indexRobot].Position.Row)
		}
	}

	return board, nil
}

func convPosToByte(pB *byte, column uint8, row uint8) {
	*pB = column<<4 + row
}
