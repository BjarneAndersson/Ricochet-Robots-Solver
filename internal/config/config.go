package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

// Mode Represents the different execution modes and corresponding output configurations
type Mode struct {
	NodeNeighbors          bool `yaml:"node_neighbors"`
	BoardStates            bool `yaml:"board_states"`
	RobotStoppingPositions bool `yaml:"robot_stopping_positions"`
}

// Config Represents the yaml configuration as struct
type Config struct {
	Mode              string                     `yaml:"mode"`
	BoardDataLocation string                     `yaml:"board_data_location"`
	Modes             map[string]map[string]Mode `yaml:"modes"`
}

// GetConfig Transform the yaml configuration into a Config struct
func GetConfig(filename string) (conf Config, err error) {
	yamlFile, err := os.ReadFile("config\\" + filename + ".yaml")
	if err != nil {
		return Config{}, err
	}

	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		return Config{}, err
	}

	return conf, nil
}
