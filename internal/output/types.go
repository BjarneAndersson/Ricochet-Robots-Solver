package output

import "Ricochet-Robot-Solver/internal/types"

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
