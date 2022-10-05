package tracker

import (
	"../config"
	"../types"
	"time"
)

func Duration(f func()) time.Duration {
	start := time.Now()

	f()

	return time.Since(start)
}

type TrackingDataSolver struct {
	InitializedBoardStates uint
	EvaluatedBoardStates   uint
	Duration               time.Duration
}

func TrackSolver(solver func(*types.GameRound, types.BoardState, *types.RobotStoppingPositions, config.Config) (TrackingDataSolver, []types.BoardState, error), gameRound *types.GameRound, initBoardState types.BoardState, robotStoppingPositions *types.RobotStoppingPositions, conf config.Config) (path []types.BoardState, trackingData TrackingDataSolver, err error) {
	start := time.Now()

	trackingData, path, err = solver(gameRound, initBoardState, robotStoppingPositions, conf)

	trackingData.Duration = time.Since(start)

	return path, trackingData, err
}
