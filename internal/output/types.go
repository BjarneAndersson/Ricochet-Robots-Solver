package output

import (
	"Ricochet-Robot-Solver/internal/input"
	"Ricochet-Robot-Solver/internal/types"
)

type OutRobotStoppingPositions []OutRobotStoppingPosition

type OutRobotStoppingPosition struct {
	Node      types.Position          `json:"node" yaml:"node"`
	Direction OutDirectionAndPosition `json:"direction" yaml:"direction"`
}

type OutDirectionAndPosition struct {
	Up    types.Position `json:"up" yaml:"up"`
	Down  types.Position `json:"down" yaml:"down"`
	Left  types.Position `json:"left" yaml:"left"`
	Right types.Position `json:"right" yaml:"right"`
}

type OutNeighbors struct {
	Node       types.Position   `json:"node" yaml:"node"`
	Directions input.Directions `json:"directions" yaml:"directions"`
}

type PathStep struct {
	Robots         input.Robots `json:"robots" yaml:"robots"`
	MovedRobot     string       `json:"moved_robot" yaml:"moved_robot"`
	MovedDirection string       `json:"moved_direction" yaml:"moved_direction"`
}
