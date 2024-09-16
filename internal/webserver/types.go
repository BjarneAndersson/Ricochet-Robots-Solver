package webserver

import "Ricochet-Robot-Solver/internal/input"

// Cell structure for rendering the board
//
// Fields:
//   - Walls: Directions of the walls in the cell.
//   - Target: Target information in the cell.
type Cell struct {
	Walls  input.Directions
	Target input.Target
}

// Board structure for organizing cells
//
// Fields:
//   - Rows: A 2D slice of cells (16x16) representing the board.
type Board struct {
	Rows [][]Cell
}
