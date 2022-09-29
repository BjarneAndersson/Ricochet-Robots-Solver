package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Mode struct {
	NodeNeighbors bool `yaml:"node_neighbors"`
	BoardStates   bool `yaml:"board_states"`
}

type Config struct {
	Mode              string          `yaml:"mode"`
	BoardDataLocation string          `yaml:"board_data_location"`
	Modes             map[string]Mode `yaml:"modes"`
}

func GetConfig(filename string) (c Config, err error) {

	yamlFile, err := os.ReadFile(filename + ".yaml")
	if err != nil {
		return Config{}, err
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}
