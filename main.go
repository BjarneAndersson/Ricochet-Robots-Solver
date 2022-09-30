package main

import (
	"./config"
	"./input"
	"./output"
	"./solver"
	"./tracker"
	"log"
)

func main() {
	// get conf
	conf, err := config.GetConfig("config")

	// convert data into board object
	board, initBoardState, robotStoppingPositions, err := input.GetData(conf.BoardDataLocation)
	if err != nil {
		log.Printf("Error loading board data:\n%v\n", err)
		return
	}

	if conf.Modes[conf.Mode]["output"].NodeNeighbors == true {
		err = output.Neighbors(&board)
		if err != nil {
			return
		}
	}

	if conf.Modes[conf.Mode]["output"].RobotStoppingPositions == true {
		err = output.RobotStoppingPositions(&robotStoppingPositions)
		if err != nil {
			return
		}
	}

	// solve the board

	path, trackingData, err := tracker.TrackSolver(solver.Solver, &board, initBoardState, &robotStoppingPositions, conf)
	if err != nil {
		log.Printf("\nError solving:\n%v\n", err)
		return
	}
	err = output.Path(path, trackingData, board.RobotColors)
	if err != nil {
		return
	}
}
