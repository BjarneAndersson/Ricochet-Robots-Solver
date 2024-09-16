package webserver

import (
	"Ricochet-Robot-Solver/internal/input"
	"encoding/json"
	"html/template"
	"net/http"
	"os"
)

// loadBoardData loads the board data from a JSON file.
//
// Parameters:
//   - filename: The name of the JSON file containing the board data.
//
// Returns:
//   - input.Board: The board data loaded from the JSON file.
//   - error: An error if one occurred while loading or parsing the file.
func loadBoardData(filename string) (input.Board, error) {
	var boardData input.Board

	// Open and read the JSON file
	file, err := os.Open(filename)
	if err != nil {
		return boardData, err
	}
	defer file.Close()

	bytes, err := os.ReadFile(filename)
	if err != nil {
		return boardData, err
	}

	// Unmarshal the JSON data into Go structs
	err = json.Unmarshal(bytes, &boardData)
	if err != nil {
		return boardData, err
	}

	return boardData, nil
}

// createBoard converts the board data to a 2D array of cells.
//
// Parameters:
//   - boardData: The board data loaded from the JSON file.
//
// Returns:
//   - Board: A 2D array of cells representing the board with walls and targets.
func createBoard(boardData input.Board) Board {
	// Create an empty 16x16 board
	board := Board{Rows: make([][]Cell, 16)}
	for i := range board.Rows {
		board.Rows[i] = make([]Cell, 16)
	}

	// Place walls in the board
	for _, wall := range boardData.Walls {
		row := wall.Position.Row
		col := wall.Position.Column
		board.Rows[row][col].Walls = wall.Directions
	}

	// Place targets in the board
	for _, target := range boardData.Targets {
		row := target.Position.Row
		col := target.Position.Column
		board.Rows[row][col].Target = target
	}

	return board
}

// renderBoard handles the rendering of the board using a template.
//
// Returns:
//   - http.HandlerFunc: A function that handles HTTP requests and writes responses.
func (s *Server) renderBoard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Load the board data from JSON
		boardData, err := loadBoardData("logs/boards/board.json")
		if err != nil {
			http.Error(w, "Failed to load board data", http.StatusInternalServerError)
			return
		}

		// Create the board with walls and targets
		board := createBoard(boardData)

		s.Logger.Info("Root route accessed")

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// Parse the template
		html, err := template.ParseFiles("website/templates/index.html")
		if err != nil {
			s.Logger.Error(err, "Error parsing template")
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}

		// Render the template
		err = html.Execute(w, board)
		if err != nil {
			s.Logger.Error(err, "Error formatting template")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
