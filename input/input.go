package input

import (
	"../helper"
	"../precomputation"
	"../types"
	"encoding/json"
	"os"
)

func GetData() (board types.Board, err error) {
	data, err := getJsonData("K:\\Coding\\Python\\Games\\Ricochet-Robots\\src\\board_data.json")
	if err != nil {
		return types.Board{}, err
	}

	board, err = loadData(data)
	if err != nil {
		return types.Board{}, err
	}

	return board, nil
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

	return board, nil
}

func convPosToByte(pB *byte, column uint8, row uint8) {
	*pB = column<<4 + row
}
