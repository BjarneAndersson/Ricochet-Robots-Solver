package input

import (
	"Ricochet-Robot-Solver/internal/config"
	"Ricochet-Robot-Solver/internal/types"
	"os"
)

// getJsonData Loads the local json data into a byte array
func getJsonData(path string) (jsonData []byte, err error) {
	jsonData, err = os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

// GetData Manages the complete data loading process and returns the formatted data
func GetData(cfg config.Config) (encodedBoard types.Board, initBoardState types.BoardState, robotOrder types.RobotOrder, err error) {
	// load the board, robots and target from corresponding json files and config
	board, err := getBoard(cfg)
	if err != nil {
		return types.Board{}, 0, 0, err
	}

	robots, err := getRobots(cfg, robotOrder)
	if err != nil {
		return types.Board{}, 0, 0, err
	}

	target, err := getTarget(cfg)
	if err != nil {
		return types.Board{}, 0, 0, err
	}

	// output the board as json
	err = board.outputJSON("./logs/boards/board.json")
	if err != nil {
		return types.Board{}, 0, 0, err
	}

	// evaluate robot order
	robotOrder = robots.getRobotOrder(target.Color)
	robots.sort(robotOrder)

	// encode game objects to solver format
	encodedBoard, err = board.encode()
	if err != nil {
		return types.Board{}, 0, 0, err
	}

	encodedTarget, err := target.encode()
	if err != nil {
		return types.Board{}, 0, 0, err

	}

	// evaluate initial robot positions
	initBoardState, err = robots.encode()
	if err != nil {
		return types.Board{}, 0, 0, err
	}

	// add target to board
	encodedBoard.Target = encodedTarget

	return encodedBoard, initBoardState, robotOrder, nil
}
