package main

import (
	"./input"
	"./output"
	"./precomputation"
	"./solver"
	"log"
)

func main() {
	// convert data into board object
	board, initBoardState, err := input.GetData()
	if err != nil {
		log.Printf("Error loading board data:\n%v\n", err)
		return
	}
	log.Printf("%+v\n", board)

	// calculate stopping positions for each node
	robotMoves, err := precomputation.PrecomputeRobotMoves(&board)
	log.Printf("\n%+v\n", robotMoves)

	// solve the board

	path, err := solver.Solver(&board, initBoardState)
	if err != nil {
		log.Printf("\nError solving:\n%v\n", err)
		return
	}
	log.Printf("%+v\n", path)
	err = output.Path(path, board.RobotColors)
	if err != nil {
		return
	}
}
