package main

import (
	"./config"
	"./input"
	"./output"
	"./precomputation"
	"./solver"
	"./tracker"
	"log"
)

func main() {
	// get conf
	conf, err := config.GetConfig("config")

	// convert data into board object
	board, initBoardState, err := input.GetData(conf.BoardDataLocation)
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

	// calculate stopping positions for each node
	robotMoves, err := precomputation.PrecomputeRobotMoves(&board)
	log.Printf("\n%+v\n", robotMoves)

	// solve the board

	path, trackingData, err := tracker.TrackSolver(solver.Solver, &board, initBoardState)
	if err != nil {
		log.Printf("\nError solving:\n%v\n", err)
		return
	}
	err = output.Path(path, trackingData, board.RobotColors)
	if err != nil {
		return
	}
}
