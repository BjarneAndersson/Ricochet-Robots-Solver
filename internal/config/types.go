package config

import "Ricochet-Robot-Solver/internal/types"

// Config structure for the application
type Config struct {
	Mode   string `yaml:"mode" env:"APPLICATION_MODE" env-description:"Mode to run the app" env-default:"deployment"`
	Server struct {
		Host string `yaml:"host" env:"SRV_HOST,HOST,HOST_NAME" env-description:"Server host" env-default:"localhost"`
		Port string `yaml:"port" env:"SRV_PORT,PORT,HOST_PORT" env-description:"Server port" env-default:"80"`
	} `yaml:"server"`
	Logging struct {
		Name        string `yaml:"name" env:"LOG_NAME" env-description:"Log name" env-default:"ricochet-robots-solver"`
		Rotation    string `yaml:"rotation" env:"LOG_ROTATION" env-description:"Log rotation" env-default:"daily"`
		Application struct {
			File struct {
				Name string `yaml:"name" env:"LOG_FILE_NAME" env-description:"Log file name" env-default:"ricochet-robots-solver"`
				Path string `yaml:"path" env:"LOG_FILE_PATH" env-description:"Log file path" env-default:"./logs"`
			} `yaml:"file"`
		} `yaml:"application"`
		Solver struct {
			File struct {
				Name string `yaml:"name" env:"LOG_SOLVER_FILE_NAME" env-description:"Log file name of solver" env-default:"datetime"`
				Path string `yaml:"path" env:"LOG_SOLVER_FILE_PATH" env-description:"Log file path of solver" env-default:"./logs/solver"`
			} `yaml:"file"`
		} `yaml:"solver"`
	} `yaml:"logging"`
	Board struct {
		Quarters struct {
			Yellow struct {
				Version  string `yaml:"version" env:"BOARD_YELLOW_VERSION" env-description:"Yellow quarter board version" env-default:"1"`
				Position string `yaml:"position" env:"BOARD_YELLOW_POSITION" env-description:"Yellow quarter board position" env-default:"NW"`
			} `yaml:"yellow"`
			Red struct {
				Version  string `yaml:"version" env:"BOARD_RED_VERSION" env-description:"Red quarter board version" env-default:"1"`
				Position string `yaml:"position" env:"BOARD_RED_POSITION" env-description:"Red quarter board position" env-default:"NE"`
			} `yaml:"red"`
			Green struct {
				Version  string `yaml:"version" env:"BOARD_GREEN_VERSION" env-description:"Green quarter board version" env-default:"1"`
				Position string `yaml:"position" env:"BOARD_GREEN_POSITION" env-description:"Green quarter board position" env-default:"SW"`
			} `yaml:"green"`
			Blue struct {
				Version  string `yaml:"version" env:"BOARD_BLUE_VERSION" env-description:"Blue quarter board version" env-default:"1"`
				Position string `yaml:"position" env:"BOARD_BLUE_POSITION" env-description:"Blue quarter board position" env-default:"SE"`
			} `yaml:"blue"`
			Location string `yaml:"location" env:"BOARD_LOCATION" env-description:"Location of the board json file" env-default:"./data/quarter-boards"`
		} `yaml:"quarters"`
		Robots struct {
			Yellow struct {
				Position types.Position `yaml:"position" env:"ROBOT_YELLOW_POSITION" env-description:"Yellow robot position" env-default:"{\"column\": 0, \"row\": 0}"`
			} `yaml:"yellow"`
			Red struct {
				Position types.Position `yaml:"position" env:"ROBOT_RED_POSITION" env-description:"Red robot position" env-default:"{\"column\": 0, \"row\": 0}"`
			} `yaml:"red"`
			Green struct {
				Position types.Position `yaml:"position" env:"ROBOT_GREEN_POSITION" env-description:"Green robot position" env-default:"{\"column\": 0, \"row\": 0}"`
			} `yaml:"green"`
			Blue struct {
				Position types.Position `yaml:"position" env:"ROBOT_BLUE_POSITION" env-description:"Blue robot position" env-default:"{\"column\": 0, \"row\": 0}"`
			} `yaml:"blue"`
		} `yaml:"robots"`
		Target struct {
			Color  string `yaml:"color" env:"TARGET_COLOR" env-description:"Target color" env-default:"yellow"`
			Symbol string `yaml:"symbol" env:"TARGET_SYMBOL" env-description:"Target symbol" env-default:"circle"`
		} `yaml:"target"`
	} `yaml:"board"`
}
