package webserver

import (
	"Ricochet-Robot-Solver/internal/input"
	"Ricochet-Robot-Solver/internal/types"
)

// Cell structure for rendering the board
//
// Fields:
//   - Walls: Directions of the walls in the cell.
//   - Target: Target information in the cell.
type Cell struct {
	Walls    input.Directions
	Target   input.Target
	Robot    input.Robot
	Position types.Position
}

// Board structure for organizing cells
//
// Fields:
//   - Rows: A 2D slice of cells (16x16) representing the board.
type Board struct {
	Rows [][]Cell
}

type ReferencesStylingAndJavascript struct {
	CSSFiles []string
	JSFiles  []string
}

type WebsiteData struct {
	CSSandJS     ReferencesStylingAndJavascript
	Board        Board
	Targets      []input.Target
	ActiveTarget input.Target
	Robots       input.Robots
}
