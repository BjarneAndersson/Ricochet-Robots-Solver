package config

import (
	// Standard library packages
	"fmt"

	// External packages
	"github.com/ilyakaznacheev/cleanenv"
	// Internal packages
)

// ReadConfig reads the configuration from the file and environment variables
func ReadConfig(config *Config) (err error) { // read configuration from the file and environment variables
	err = cleanenv.ReadConfig("./config/config.yml", config)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}
