package main

import (
	// Standard library packages
	"fmt"

	// External packages

	// Internal packages
	config "Ricochet-Robot-Solver/internal/config"
	"Ricochet-Robot-Solver/internal/input"
	logging "Ricochet-Robot-Solver/internal/log"
	"Ricochet-Robot-Solver/internal/output"
	solver "Ricochet-Robot-Solver/internal/solver"
	"Ricochet-Robot-Solver/internal/tracker"
)

func main() {
	var err error
	var cfg config.Config
	var logger logging.Logger
	var solveRecorder logging.SolveRecorder

	// CONFIG
	err = config.ReadConfig(&cfg)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	// LOGGER
	logger, solveRecorder = logging.Setup(cfg)

	// APPLICATION
	logger.Info("Starting Application")

	// transform json data to board object
	board, initBoardState, initRobotOrder, robotStoppingPositions, err := input.GetData(cfg.BoardDataLocation)
	if err != nil {
		logger.Fatal(err, "Loading board data")
		return
	}

	// output extra information based on config
	err = output.Neighbors(&board, &solveRecorder)
	if err != nil {
		logger.Warn(err, "Outputting neighbors")
	}

	err = output.RobotStoppingPositions(&robotStoppingPositions, &solveRecorder)
	if err != nil {
		logger.Warn(err, "Outputting robot stopping positions")
	}

	// solve the board (with tracking)
	path, trackingData, err := tracker.TrackSolver(solver.Solver, &board, initBoardState, &robotStoppingPositions, cfg)
	if err != nil {
		logger.Error(err, "Solving")
		return
	}

	// output the path
	err = output.Path(path, trackingData, initRobotOrder, &solveRecorder)
	if err != nil {
		logger.Fatal(err, "Outputting path")
		return
	}
}
