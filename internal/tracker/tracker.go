package tracker

import (
	"Ricochet-Robot-Solver/internal/config"
	"Ricochet-Robot-Solver/internal/types"
	"time"
)

// TrackingDataSolver Collects the tracking data in a central structure
type TrackingDataSolver struct {
	InitializedBoardStates uint          `json:"InitializedBoardStates"`
	EvaluatedBoardStates   uint          `json:"EvaluatedBoardStates"`
	Duration               time.Duration `json:"Duration"`
}

// TrackSolver Tracks the solver during the execution and returns the measurements
func TrackSolver(solver func(*types.Board, types.BoardState, *types.RobotStoppingPositions, config.Config) (TrackingDataSolver, []types.BoardState, error), board *types.Board, initBoardState types.BoardState, robotStoppingPositions *types.RobotStoppingPositions, conf config.Config) (path []types.BoardState, trackingData TrackingDataSolver, err error) {
	start := time.Now()

	trackingData, path, err = solver(board, initBoardState, robotStoppingPositions, conf)
	if err != nil {
		return path, trackingData, err
	}

	trackingData.Duration = time.Since(start)

	if trackingData.Duration == 0 { // Prevent division by zero
		trackingData.Duration = 1
	}

	return path, trackingData, err
}
