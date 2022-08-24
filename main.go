package main

import (
	"./input"
	"./solver"
	"fmt"
)

func main() {
	board, initBoardState, err := input.GetData()
	if err != nil {
		fmt.Printf("Error loading board data:\n%v\n", err)
	}

	fmt.Printf("%+v\n", board)

	solver.Solver(&board, initBoardState)
}
