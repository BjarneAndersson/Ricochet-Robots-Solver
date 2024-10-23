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

func createTargets() (targets []input.Target) {
	// Define the target data
	targetsData := []struct {
		Symbol string
		Color  string
	}{
		{Symbol: "circle", Color: "yellow"},
		{Symbol: "triangle", Color: "yellow"},
		{Symbol: "square", Color: "yellow"},
		{Symbol: "hexagon", Color: "yellow"},
		{Symbol: "circle", Color: "red"},
		{Symbol: "triangle", Color: "red"},
		{Symbol: "square", Color: "red"},
		{Symbol: "hexagon", Color: "red"},
		{Symbol: "circle", Color: "green"},
		{Symbol: "triangle", Color: "green"},
		{Symbol: "square", Color: "green"},
		{Symbol: "hexagon", Color: "green"},
		{Symbol: "circle", Color: "blue"},
		{Symbol: "triangle", Color: "blue"},
		{Symbol: "square", Color: "blue"},
		{Symbol: "hexagon", Color: "blue"},
		{Symbol: "vertex", Color: "multicolor"},
	}

	// Create the targets
	targets = make([]input.Target, len(targetsData))
	for i, target := range targetsData {
		targets[i] = input.Target{
			Color: target.Color,
			Shape: target.Symbol,
		}
	}

	return targets
}

// renderBoard handles the rendering of the board using a template.
//
// Returns:
//   - http.HandlerFunc: A function that handles HTTP requests and writes responses.
func (s *Server) renderBoard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		// load the board
		err = s.loadBoard()
		if err != nil {
			s.Logger.Error(err, "Loading board data")
			http.Error(w, "Failed to load board data", http.StatusInternalServerError)
			return
		}

		// Load the board data from JSON
		boardData, err := loadBoardData("logs/boards/board.json")
		if err != nil {
			http.Error(w, "Failed to load board data", http.StatusInternalServerError)
			return
		}

		// Create the board with walls and targets
		board := createBoard(boardData)
		targets := createTargets()
		activeTarget := input.Target{
			Color: s.Config.Board.Target.Color,
			Shape: s.Config.Board.Target.Symbol,
		}

		// Create the website data
		websiteData := WebsiteData{
			Board:        board,
			Targets:      targets,
			ActiveTarget: activeTarget,
		}

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
		err = html.Execute(w, websiteData)
		if err != nil {
			s.Logger.Error(err, "Error formatting template")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
