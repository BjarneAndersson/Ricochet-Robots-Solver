package tracker

import (
	"../config"
	"../types"
	"time"
)

type TrackingDataSolver struct {
	InitializedBoardStates uint
	EvaluatedBoardStates   uint
	Duration               time.Duration
}

func TrackSolver(solver func(*types.Board, types.BoardState, *types.RobotStoppingPositions, config.Config) (TrackingDataSolver, []types.BoardState, error), board *types.Board, initBoardState types.BoardState, robotStoppingPositions *types.RobotStoppingPositions, conf config.Config) (path []types.BoardState, trackingData TrackingDataSolver, err error) {
	start := time.Now()

	trackingData, path, err = solver(board, initBoardState, robotStoppingPositions, conf)

	trackingData.Duration = time.Since(start)

	return path, trackingData, err
}
