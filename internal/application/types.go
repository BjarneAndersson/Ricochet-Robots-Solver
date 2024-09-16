package application

import (
	"Ricochet-Robot-Solver/internal/tracker"
	"Ricochet-Robot-Solver/internal/types"
)

// AppState holds the state of the application.
//
// Fields:
//   - Board: The current state of the game board.
//   - InitBoardState: The initial state of the game board.
//   - InitRobotOrder: The initial order of the robots.
//   - RobotStoppingPositions: The positions where robots stop.
//   - Path: The sequence of board states representing the path taken.
//   - TrackingData: Data used for tracking the solver's progress.
type AppState struct {
	Board                  types.Board
	InitBoardState         types.BoardState
	InitRobotOrder         types.RobotOrder
	RobotStoppingPositions types.RobotStoppingPositions
	Path                   []types.BoardState
	TrackingData           tracker.TrackingDataSolver
}
