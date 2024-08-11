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
	// define the quarter boards
	var quarterBoards struct {
		yellow QuarterBoard
		red    QuarterBoard
		green  QuarterBoard
		blue   QuarterBoard
	}

	// check: every position is present
	set := make(map[string]struct{})
	set[cfg.Board.Quarters.Yellow.Position] = struct{}{}
	set[cfg.Board.Quarters.Red.Position] = struct{}{}
	set[cfg.Board.Quarters.Green.Position] = struct{}{}
	set[cfg.Board.Quarters.Blue.Position] = struct{}{}
	if len(set) < 4 {
		return Board{}, fmt.Errorf("not all quarter boards are present")
	}

	// load the quarter boards from the json files
	err = quarterBoards.yellow.loadFileData(cfg.Board.Quarters.Location, "yellow", cfg.Board.Quarters.Yellow.Version)
	if err != nil {
		return Board{}, err
	}
	err = quarterBoards.red.loadFileData(cfg.Board.Quarters.Location, "red", cfg.Board.Quarters.Red.Version)
	if err != nil {
		return Board{}, err
	}
	err = quarterBoards.green.loadFileData(cfg.Board.Quarters.Location, "green", cfg.Board.Quarters.Green.Version)
	if err != nil {
		return Board{}, err
	}
	err = quarterBoards.blue.loadFileData(cfg.Board.Quarters.Location, "blue", cfg.Board.Quarters.Blue.Version)
	if err != nil {
		return Board{}, err
	}

	// define position of the quarter board on the board
	quarterBoards.yellow.Position = cfg.Board.Quarters.Yellow.Position
	quarterBoards.red.Position = cfg.Board.Quarters.Red.Position
	quarterBoards.green.Position = cfg.Board.Quarters.Green.Position
	quarterBoards.blue.Position = cfg.Board.Quarters.Blue.Position

	// adjust the rotation of the quarter boards
	quarterBoards.yellow.adjustRotation()
	quarterBoards.red.adjustRotation()
	quarterBoards.green.adjustRotation()
	quarterBoards.blue.adjustRotation()

	// stitch the quarter boards together
	board, err = stitchBoard([4]QuarterBoard{quarterBoards.yellow, quarterBoards.red, quarterBoards.green, quarterBoards.blue})
	if err != nil {
		return Board{}, err
	}

	return board, nil
}

// stitchBoard Stitches the four quarter boards together to create the complete board
func stitchBoard(quarterBoards [4]QuarterBoard) (board Board, err error) {
	// create the board
	for _, quarterBoard := range quarterBoards {
		var columnOffset byte = 0
		var rowOffset byte = 0

		switch quarterBoard.Position {
		case "NW":
			columnOffset = 0
			rowOffset = 0
		case "NE":
			columnOffset = 8
			rowOffset = 0
		case "SW":
			columnOffset = 0
			rowOffset = 8
		case "SE":
			columnOffset = 8
			rowOffset = 8
		}

		// add walls with corresponding offset to board
		for _, wall := range quarterBoard.Walls {
			board.Walls = append(board.Walls, Wall{
				Position: types.Position{
					Column: wall.Position.Column + columnOffset,
					Row:    wall.Position.Row + rowOffset,
				},
				Directions: Directions{
					Top:    wall.Directions.Top,
					Bottom: wall.Directions.Bottom,
					Left:   wall.Directions.Left,
					Right:  wall.Directions.Right,
				},
			})
		}

		// add targets with corresponding offset to board
		for _, target := range quarterBoard.Targets {
			board.Targets = append(board.Targets, Target{
				Position: types.Position{
					Column: target.Position.Column + columnOffset,
					Row:    target.Position.Row + rowOffset,
				},
				Color: target.Color,
				Shape: target.Shape,
			})
		}
	}

	// add the walls at the stitch lines to both sides
	for _, wall := range board.Walls {
		if wall.Position.Column == 7 && wall.Directions.Right == true {
			board.Walls = append(board.Walls, Wall{
				Position: types.Position{
					Column: 8,
					Row:    wall.Position.Row,
				},
				Directions: Directions{
					Top:    false,
					Bottom: false,
					Left:   true,
					Right:  false,
				},
			})
		} else if wall.Position.Column == 8 && wall.Directions.Left == true {
			board.Walls = append(board.Walls, Wall{
				Position: types.Position{
					Column: 7,
					Row:    wall.Position.Row,
				},
				Directions: Directions{
					Top:    false,
					Bottom: false,
					Left:   false,
					Right:  true,
				},
			})
		}

		if wall.Position.Row == 7 && wall.Directions.Bottom == true {
			board.Walls = append(board.Walls, Wall{
				Position: types.Position{
					Column: wall.Position.Column,
					Row:    8,
				},
				Directions: Directions{
					Top:    true,
					Bottom: false,
					Left:   false,
					Right:  false,
				},
			})
		} else if wall.Position.Row == 8 && wall.Directions.Top == true {
			board.Walls = append(board.Walls, Wall{
				Position: types.Position{
					Column: wall.Position.Column,
					Row:    7,
				},
				Directions: Directions{
					Top:    false,
					Bottom: true,
					Left:   false,
					Right:  false,
				},
			})
		}
	}

	return board, nil
}

func (board Board) encode() (encodedBoard types.Board, err error) {
	// node conversion
	for rowIndex := 0; rowIndex < 16; rowIndex++ {
		for columnIndex := 0; columnIndex < 16; columnIndex++ {
			var node types.Node

			// distance between node and target
			helper.SetBit((*byte)(&node), 7, true)
			helper.SetBit((*byte)(&node), 6, true)
			helper.SetBit((*byte)(&node), 5, true)

			// is a robot present
			helper.SetBit((*byte)(&node), 4, false)

			// neighbors
			helper.SetBit((*byte)(&node), 3, true)
			helper.SetBit((*byte)(&node), 2, true)
			helper.SetBit((*byte)(&node), 1, true)
			helper.SetBit((*byte)(&node), 0, true)

			encodedBoard.Grid[rowIndex][columnIndex] = node
		}
	}

	// add walls to the board -> set the designated  bits of the nodes
	for _, wall := range board.Walls {
		node := &encodedBoard.Grid[wall.Position.Row][wall.Position.Column]

		if wall.Directions.Top {
			helper.SetBit((*byte)(node), 3, false)
		}
		if wall.Directions.Bottom {
			helper.SetBit((*byte)(node), 2, false)
		}
		if wall.Directions.Left {
			helper.SetBit((*byte)(node), 1, false)
		}
		if wall.Directions.Right {
			helper.SetBit((*byte)(node), 0, false)
		}
	}

	return encodedBoard, err
}

func (board Board) outputJSON(filePath string) error {
	// Convert the board struct to JSON
	boardJSON, err := json.Marshal(board)
	if err != nil {
		fmt.Println("Error marshalling board to JSON:", err)
		return err
	}

	// Open a new log file, overwriting any existing file
	fBoard, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer func(fBoard *os.File) {
		err := fBoard.Close()
		if err != nil {
			panic(err)
		}
	}(fBoard)

	// Write the JSON content to the file
	_, err = fBoard.Write(boardJSON)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return err
	}

	return nil
}
