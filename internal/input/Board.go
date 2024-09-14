package input

import (
	"Ricochet-Robot-Solver/internal/config"
	"Ricochet-Robot-Solver/internal/helper"
	"Ricochet-Robot-Solver/internal/types"
	"encoding/json"
	"fmt"
	"os"
)

func getBoard(cfg config.Config) (board Board, err error) {
	// Define the quarter boards
	var quarterBoards struct {
		yellow, red, green, blue QuarterBoard
	}

	// Check if every position is present
	positions := map[string]struct{}{
		cfg.Board.Quarters.Yellow.Position: {},
		cfg.Board.Quarters.Red.Position:    {},
		cfg.Board.Quarters.Green.Position:  {},
		cfg.Board.Quarters.Blue.Position:   {},
	}
	if len(positions) < 4 {
		return Board{}, fmt.Errorf("not all quarter boards are present")
	}

	// Load the quarter boards from the JSON files
	for color, qb := range map[string]*QuarterBoard{
		"yellow": &quarterBoards.yellow,
		"red":    &quarterBoards.red,
		"green":  &quarterBoards.green,
		"blue":   &quarterBoards.blue,
	} {
		var version string
		switch color {
		case "yellow":
			version = cfg.Board.Quarters.Yellow.Version
		case "red":
			version = cfg.Board.Quarters.Red.Version
		case "green":
			version = cfg.Board.Quarters.Green.Version
		case "blue":
			version = cfg.Board.Quarters.Blue.Version
		}
		if err := qb.loadFileData(cfg.Board.Quarters.Location, color, version); err != nil {
			return Board{}, err
		}
	}

	// Define position and adjust rotation of the quarter boards
	for color, qb := range map[string]*QuarterBoard{
		"yellow": &quarterBoards.yellow,
		"red":    &quarterBoards.red,
		"green":  &quarterBoards.green,
		"blue":   &quarterBoards.blue,
	} {
		switch color {
		case "yellow":
			qb.Position = cfg.Board.Quarters.Yellow.Position
		case "red":
			qb.Position = cfg.Board.Quarters.Red.Position
		case "green":
			qb.Position = cfg.Board.Quarters.Green.Position
		case "blue":
			qb.Position = cfg.Board.Quarters.Blue.Position
		}
		qb.adjustRotation()
	}

	// Stitch the quarter boards together
	return stitchBoard([4]QuarterBoard{quarterBoards.yellow, quarterBoards.red, quarterBoards.green, quarterBoards.blue})
}

// stitchBoard Stitches the four quarter boards together to create the complete board
func stitchBoard(quarterBoards [4]QuarterBoard) (board Board, err error) {
	// Define the offsets for each quarter board position
	offsets := map[string]types.Position{
		"NW": {0, 0}, "NE": {8, 0}, "SW": {0, 8}, "SE": {8, 8},
	}

	// Iterate over each quarter board
	for _, qb := range quarterBoards {
		offset := offsets[qb.Position]

		// Adjust the positions of walls and add them to the board
		for _, wall := range qb.Walls {
			board.Walls = append(board.Walls, Wall{
				Position: types.Position{
					Column: wall.Position.Column + offset.Column,
					Row:    wall.Position.Row + offset.Row,
				},
				Directions: wall.Directions,
			})
		}

		// Adjust the positions of targets and add them to the board
		for _, target := range qb.Targets {
			board.Targets = append(board.Targets, Target{
				Position: types.Position{
					Column: target.Position.Column + offset.Column,
					Row:    target.Position.Row + offset.Row,
				},
				Color: target.Color,
				Shape: target.Shape,
			})
		}
	}

	// Handle the stitching of walls at the boundaries between quarter boards
	for _, wall := range board.Walls {
		switch {
		case wall.Position.Column == 7 && wall.Directions.Right:
			board.Walls = append(board.Walls, Wall{
				Position:   types.Position{Column: 8, Row: wall.Position.Row},
				Directions: Directions{Left: true},
			})
		case wall.Position.Column == 8 && wall.Directions.Left:
			board.Walls = append(board.Walls, Wall{
				Position:   types.Position{Column: 7, Row: wall.Position.Row},
				Directions: Directions{Right: true},
			})
		case wall.Position.Row == 7 && wall.Directions.Bottom:
			board.Walls = append(board.Walls, Wall{
				Position:   types.Position{Column: wall.Position.Column, Row: 8},
				Directions: Directions{Top: true},
			})
		case wall.Position.Row == 8 && wall.Directions.Top:
			board.Walls = append(board.Walls, Wall{
				Position:   types.Position{Column: wall.Position.Column, Row: 7},
				Directions: Directions{Bottom: true},
			})
		}
	}

	return board, nil
}

// encode: Encodes the board to the solver format
func (board Board) encode() (encodedBoard types.Board, err error) {
	// node conversion
	for rowIndex := 0; rowIndex < 16; rowIndex++ {
		for columnIndex := 0; columnIndex < 16; columnIndex++ {
			// bit 7-5: distance between node and target
			// bit 4: is a robot present
			// bit 3-0: neighbors of the node
			encodedBoard.Grid[rowIndex][columnIndex] = types.Node(0xE0) // 0xE0 sets bits 7, 6, 5, 3, 2, 1, 0 to true

		}
	}

	// add walls to the board -> set the designated neighbor bits of the nodes to false
	for _, wall := range board.Walls {
		node := &encodedBoard.Grid[wall.Position.Row][wall.Position.Column]

		if wall.Directions.Top { // top | set bit 7 to false
			helper.SetBit((*byte)(node), 3, false)
		}
		if wall.Directions.Bottom { // bottom | set bit 6 to false
			helper.SetBit((*byte)(node), 2, false)
		}
		if wall.Directions.Left { // left | set bit 1 to false
			helper.SetBit((*byte)(node), 1, false)
		}
		if wall.Directions.Right { // right | set bit 0 to false
			helper.SetBit((*byte)(node), 0, false)
		}
	}

	return encodedBoard, err
}

// outputJSON: Outputs the board as JSON to the specified file path
func (board Board) outputJSON(filePath string) error {
	// Convert the board struct to JSON
	boardJSON, err := json.Marshal(board)
	if err != nil {
		return fmt.Errorf("error marshalling board to JSON: %w", err)
	}

	// Open a new log file, overwriting any existing file
	fBoard, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer func(fBoard *os.File) {
		err := fBoard.Close()
		if err != nil {
			panic(err)
		}
	}(fBoard)

	// Write the JSON content to the file
	if _, err = fBoard.Write(boardJSON); err != nil {
		return fmt.Errorf("error writing JSON to file: %w", err)
	}

	return nil
}
