package config

// Config structure for the application
type Config struct {
	Mode              string `yaml:"mode" env:"APPLICATION_MODE" env-description:"Mode to run the app" env-default:"deployment"`
	BoardDataLocation string `yaml:"boardDataLocation" env:"BOARD_DATA_LOCATION" env-description:"Location of the board data" env-default:"test/data/9c.json"`
	Server            struct {
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
}
