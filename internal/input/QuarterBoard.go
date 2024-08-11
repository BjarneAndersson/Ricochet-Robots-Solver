package input

import (
	"Ricochet-Robot-Solver/internal/types"
	"encoding/json"
	"fmt"
)

func (quarterBoard *QuarterBoard) loadFileData(folderPath string, color string, version string) (err error) {
	// get json data from file
	jsonData, err := getJsonData(fmt.Sprintf("%s/%s_%s.json", folderPath, color, version))
	if err != nil {
		return err
	}

	// decode json data to quarter board
	err = json.Unmarshal(jsonData, quarterBoard)
	if err != nil {
		return err
	}

	return err
}

func (quarterBoard *QuarterBoard) adjustRotation() {
	// rotate the quarter board based on its position
	var rotationalSteps byte
	switch quarterBoard.Position {
	case "NW":
		rotationalSteps = RotationOfNW
	case "NE":
		rotationalSteps = RotationOfNE
	case "SW":
		rotationalSteps = RotationOfSW
	case "SE":
		rotationalSteps = RotationOfSE
	default:
		rotationalSteps = 0
	}
	quarterBoard.rotate(int(rotationalSteps))
}

// rotate Rotates the quarter board by 0, 90, 180 or 270 degrees in the right direction
func (quarterBoard *QuarterBoard) rotate(steps int) {
	const BoardSize = 8
	var rotatedWalls []Wall
	var rotatedTargets []Target

	for i := 0; i < steps; i++ {
		// rotational steps
		for i := 0; i < steps; i++ {
			// rotate the walls
			rotatedWalls = make([]Wall, len(quarterBoard.Walls))
			for i, wall := range quarterBoard.Walls {
				rotatedWalls[i] = Wall{
					Position:   rotatePosition(wall.Position, BoardSize),
					Directions: rotateDirections(wall.Directions),
				}
			}

			rotatedTargets = make([]Target, len(quarterBoard.Targets))
			for i, target := range quarterBoard.Targets {
				rotatedTargets[i] = Target{
					Position: rotatePosition(target.Position, BoardSize),
					Color:    target.Color,
					Shape:    target.Shape,
				}
			}
		}

		quarterBoard.Walls = rotatedWalls
		quarterBoard.Targets = rotatedTargets
	}
}

func rotatePosition(position types.Position, n byte) types.Position {
	return types.Position{
		Column: (n - 1) - position.Row,
		Row:    position.Column,
	}
}

func rotateDirections(directions Directions) Directions {
	return Directions{
		Top:    directions.Left,
		Bottom: directions.Right,
		Left:   directions.Bottom,
		Right:  directions.Top,
	}
}
