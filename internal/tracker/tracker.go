package tracker

import (
	"Ricochet-Robot-Solver/internal/config"
	"Ricochet-Robot-Solver/internal/types"
	"time"
)

// TrackingDataSolver Collects the tracking data in a central structure
type TrackingDataSolver struct {
	InitializedBoardStates uint
	EvaluatedBoardStates   uint
	Duration               time.Duration
}

// TrackSolver Tracks the solver during the execution and returns the measurements
func TrackSolver(solver func(*types.Board, types.BoardState, *types.RobotStoppingPositions, config.Config) (TrackingDataSolver, []types.BoardState, error), board *types.Board, initBoardState types.BoardState, robotStoppingPositions *types.RobotStoppingPositions, conf config.Config) (path []types.BoardState, trackingData TrackingDataSolver, err error) {
	start := time.Now()

	trackingData, path, err = solver(board, initBoardState, robotStoppingPositions, conf)

	trackingData.Duration = time.Since(start)

	return path, trackingData, err
}
