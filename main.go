package main

import (
	"./input"
	"./precomputation"
	"./solver"
	"fmt"
)

func main() {
	// convert data into board object
	board, initBoardState, err := input.GetData()
	if err != nil {
		fmt.Printf("Error loading board data:\n%v\n", err)
		return
	}
	fmt.Printf("%+v\n", board)

	// calculate stopping positions for each node
	robotMoves, err := precomputation.PrecomputeRobotMoves(&board)
	fmt.Printf("\n%+v\n", robotMoves)

	// solve the board

	path, err := solver.Solver(&board, initBoardState)
	if err != nil {
		fmt.Printf("\nError solving:\n%v\n", err)
		return
	}
	fmt.Printf("\n%+v\n", path)
}
