package input

import (
	"Ricochet-Robot-Solver/internal/types"
)

// RotationOfNW, RotationOfNE, RotationOfSW, RotationOfSE The quarter boards are initially all in the same orientation, through rotation the placement of the quarter boards
// for the representation of the complete board can be adjusted.
// Each constant represents the necessary rotations of the quarter board needed for each position with n*90 degrees.
const (
	RotationOfNW byte = iota // no rotation needed
	RotationOfNE byte = iota // rotation of 90 degrees needed
	RotationOfSE byte = iota // rotation of 180 degrees needed
	RotationOfSW byte = iota // rotation of 270 degrees needed
)

// Directions Used in data import to easily manage the attributes
// Boolean values to describe the presence of a wall or neighbor in a certain direction
type Directions struct {
	Top    bool `json:"top"`
	Bottom bool `json:"bottom"`
	Left   bool `json:"left"`
	Right  bool `json:"right"`
}

// Wall Used in data import to easily manage the attributes
type Wall struct {
	Position   types.Position `json:"position"`   // position of the node that walls surround
	Directions Directions     `json:"directions"` // state in which direction a wall is present
}

// Robot Used in data import to easily manage the attributes
type Robot struct {
	Color    string         `json:"color"`    // color of the robot
	Position types.Position `json:"position"` // position of the robot on the board
}

type Robots [4]Robot // slice representing the four colored robots

// Target Used in data import to easily manage the attributes
type Target struct {
	Color    string         `json:"color"`    // color of the target
	Shape    string         `json:"shape"`    // shape of the target
	Position types.Position `json:"position"` // position of the target on the board
}

// QuarterBoard Used in data import to easily manage the attributes
type QuarterBoard struct {
	Walls    []Wall   `json:"walls"`   // describes all walls that are present on the quarter board
	Targets  []Target `json:"targets"` // describes all targets that are present on the quarter board
	Position string   // Describes in which corner the quarter board is placed | [NW, NE, SW, SE]
}

// Board Used in data import to easily manage the attributes
type Board struct {
	Walls   []Wall   `json:"walls"`  // describes all walls that are present on the board
	Targets []Target `json:"target"` // describes all targets that are present on the board
}
