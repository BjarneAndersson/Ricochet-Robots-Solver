package input

import (
	"Ricochet-Robot-Solver/internal/config"
	"Ricochet-Robot-Solver/internal/helper"
	"Ricochet-Robot-Solver/internal/types"
)

// getRobots Loads the robots from the configuration file
func getRobots(cfg config.Config, robotOrder types.RobotOrder) (robots Robots, err error) {

	// define the robots
	var robotsStruct struct {
		yellow Robot
		red    Robot
		green  Robot
		blue   Robot
	}

	robotsStruct.yellow = Robot{
		Color: "yellow",
		Position: types.Position{
			Column: cfg.Board.Robots.Yellow.Position.Column,
			Row:    cfg.Board.Robots.Yellow.Position.Row,
		},
	}
	robotsStruct.red = Robot{
		Color: "red",
		Position: types.Position{
			Column: cfg.Board.Robots.Red.Position.Column,
			Row:    cfg.Board.Robots.Red.Position.Row,
		},
	}
	robotsStruct.green = Robot{
		Color: "green",
		Position: types.Position{
			Column: cfg.Board.Robots.Green.Position.Column,
			Row:    cfg.Board.Robots.Green.Position.Row,
		},
	}
	robotsStruct.blue = Robot{
		Color: "blue",
		Position: types.Position{
			Column: cfg.Board.Robots.Blue.Position.Column,
			Row:    cfg.Board.Robots.Blue.Position.Row,
		},
	}

	// convert the robots to an array
	robots = [4]Robot{robotsStruct.yellow, robotsStruct.red, robotsStruct.green, robotsStruct.blue}

	return robots, nil
}

func (robots *Robots) encode() (encodedRobotPositions types.BoardState, err error) {

	encodedRobotPositions = types.BoardState(
		uint32(helper.ConvPositionToByte(robots[0].Position))<<24 |
			uint32(helper.ConvPositionToByte(robots[1].Position))<<16 |
			uint32(helper.ConvPositionToByte(robots[2].Position))<<8 |
			uint32(helper.ConvPositionToByte(robots[3].Position))<<0)

	return encodedRobotPositions, err
}

func (robots *Robots) getRobotOrder(targetColor string) (robotOrder types.RobotOrder) {
	var colors = []string{"yellow", "red", "green", "blue"}

	switch targetColor { // robot with target color is always first
	case "yellow":
		colors = append([]string{"yellow"}, colors[1:]...)
	case "red":
		colors = append([]string{"red"}, append(colors[:1], colors[2:]...)...)
	case "green":
		colors = append([]string{"green"}, append(colors[:2], colors[3:]...)...)
	case "blue":
		colors = append([]string{"blue"}, colors[:3]...)
	}

	helper.SetRobotColorByIndex(&robotOrder, targetColor, 0)
	helper.SetRobotColorByIndex(&robotOrder, colors[1], 1)
	helper.SetRobotColorByIndex(&robotOrder, colors[2], 2)
	helper.SetRobotColorByIndex(&robotOrder, colors[3], 3)

	return robotOrder
}

func (robots *Robots) sort(robotOrder types.RobotOrder) {
	var robotsSlice [4]Robot
	robotsSlice = [4]Robot{robots[0], robots[1], robots[2], robots[3]}

	// convert the robots to an array according to the robot order
	for i := 0; i < 4; i++ {
		switch helper.GetRobotColorCodeByIndex(robotOrder, uint8(i)) {
		case types.RobotColorYellow:
			robotsSlice[i] = robots[types.RobotColorYellow]
		case types.RobotColorRed:
			robotsSlice[i] = robots[types.RobotColorRed]
		case types.RobotColorGreen:
			robotsSlice[i] = robots[types.RobotColorGreen]
		case types.RobotColorBlue:
			robotsSlice[i] = robots[types.RobotColorBlue]
		}
	}

	*robots = [4]Robot{robotsSlice[0], robotsSlice[1], robotsSlice[2], robotsSlice[3]}
}
