package main

import (
	"Ricochet-Robot-Solver/internal/config"
	"Ricochet-Robot-Solver/internal/input"
	"Ricochet-Robot-Solver/internal/output"
	"Ricochet-Robot-Solver/internal/solver"
	"Ricochet-Robot-Solver/internal/tracker"
	"log"
)

func main() {
	// get program config
	conf, err := config.GetConfig("config")

	// transform json data to board object
	board, initBoardState, initRobotOrder, robotStoppingPositions, err := input.GetData(conf.BoardDataLocation)
	if err != nil {
		log.Printf("Error loading board data:\n%v\n", err)
		return
	}

	// output extra information based on config
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

	// solve the board (with tracking)
	path, trackingData, err := tracker.TrackSolver(solver.Solver, &board, initBoardState, &robotStoppingPositions, conf)
	if err != nil {
		log.Printf("\nError solving:\n%v\n", err)
		return
	}

	// output the path
	err = output.Path(path, trackingData, initRobotOrder)
	if err != nil {
		return
	}
}
