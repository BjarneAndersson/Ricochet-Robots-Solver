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

	path, err := solver.Solver(&board, initBoardState)
	if err != nil {
		return
	}

	fmt.Printf("%+v\n", path)
}
