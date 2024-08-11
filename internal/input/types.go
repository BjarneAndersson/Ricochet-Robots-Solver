package input

import (
	"Ricochet-Robot-Solver/internal/types"
)

const ( // corner which represents the center of the board is initially placed in the NW corner
	RotationOfNW byte = iota // specifies the rotation of the quarter board for each position with n*90 degrees
	RotationOfNE byte = iota
	RotationOfSE byte = iota
	RotationOfSW byte = iota
)

// Directions Used in data import to easily manage the attributes
type Directions struct {
	Top    bool `json:"top"`
	Bottom bool `json:"bottom"`
	Left   bool `json:"left"`
	Right  bool `json:"right"`
}

// Wall Used in data import to easily manage the attributes
type Wall struct {
	Position   types.Position `json:"position"`
	Directions Directions     `json:"directions"` // state in which direction a wall is present
}

// Robot Used in data import to easily manage the attributes
type Robot struct {
	Color    string         `json:"color"`
	Position types.Position `json:"position"`
}

type Robots [4]Robot

// Target Used in data import to easily manage the attributes
type Target struct {
	Color    string         `json:"color"`
	Shape    string         `json:"shape"`
	Position types.Position `json:"position"`
}

// QuarterBoard Used in data import to easily manage the attributes
type QuarterBoard struct {
	Walls    []Wall   `json:"walls"`
	Targets  []Target `json:"targets"`
	Position string   // NW, NE, SW, SE
}

// Board Used in data import to easily manage the attributes
type Board struct {
	Walls   []Wall   `json:"walls"`
	Targets []Target `json:"target"`
}
