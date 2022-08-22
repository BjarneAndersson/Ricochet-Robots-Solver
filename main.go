package main

import (
	"./input"
	"fmt"
)

func main() {
	board, err := input.GetData()
	if err != nil {
		fmt.Printf("Error loading board data:\n%v\n", err)
	}

	fmt.Printf("%+v\n", board)
}
