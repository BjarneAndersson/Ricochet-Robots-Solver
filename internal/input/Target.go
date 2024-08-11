package input

import (
	"Ricochet-Robot-Solver/internal/config"
	"Ricochet-Robot-Solver/internal/helper"
	"Ricochet-Robot-Solver/internal/types"
)

// getTarget Loads the target from the configuration file
func getTarget(cfg config.Config) (target Target, err error) {
	target = Target{
		Position: types.Position{
			Column: cfg.Board.Target.Position.Column,
			Row:    cfg.Board.Target.Position.Row,
		},
		Color: cfg.Board.Target.Color,
		Shape: cfg.Board.Target.Symbol,
	}
	return target, nil
}

func (target Target) encode() (encodedTarget types.Target, err error) {
	var targetColorAndSymbol byte
	targetPosition := helper.ConvPositionToByte(target.Position)

	switch target.Color {
	case "yellow":
		helper.SetBit(&targetColorAndSymbol, 7, true)
	case "red":
		helper.SetBit(&targetColorAndSymbol, 6, true)
	case "green":
		helper.SetBit(&targetColorAndSymbol, 5, true)
	case "blue":
		helper.SetBit(&targetColorAndSymbol, 4, true)
	}

	switch target.Shape {
	case "circle":
		helper.SetBit(&targetColorAndSymbol, 3, true)
	case "triangle":
		helper.SetBit(&targetColorAndSymbol, 2, true)
	case "square":
		helper.SetBit(&targetColorAndSymbol, 1, true)
	case "hexagon":
		helper.SetBit(&targetColorAndSymbol, 0, true)
	case "vortex":
		helper.SetBit(&targetColorAndSymbol, 3, true)
		helper.SetBit(&targetColorAndSymbol, 2, true)
		helper.SetBit(&targetColorAndSymbol, 1, true)
		helper.SetBit(&targetColorAndSymbol, 0, true)
	}

	encodedTarget = types.Target(uint16(targetColorAndSymbol)<<8 | uint16(targetPosition))

	return encodedTarget, err
}
