package tracker

import (
	"../types"
	"time"
)

func Duration(f func()) time.Duration {
	start := time.Now()

	f()

	return time.Since(start)
}

func DurationSolver(solver func(*types.Board, types.BoardState) ([]types.BoardState, error), board *types.Board, initBoardState types.BoardState) (path []types.BoardState, duration time.Duration, err error) {
	start := time.Now()

	path, err = solver(board, initBoardState)

	return path, time.Since(start), err
}
