package main

import (
	"Ricochet-Robot-Solver/internal/webserver"
	"net/http"

	// Standard library packages
	"fmt"

	// External packages

	// Internal packages
	config "Ricochet-Robot-Solver/internal/config"
	logging "Ricochet-Robot-Solver/internal/log"
)

func main() {
	var err error
	var cfg config.Config
	var logger logging.Logger
	var solveRecorder logging.SolveRecorder

	// CONFIG
	err = config.ReadConfig(&cfg)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	// LOGGER
	logger, solveRecorder = logging.Setup(cfg)

	// APPLICATION
	logger.Info("Starting Application")

	// WEB SERVER
	server := webserver.NewServer(cfg.Server.Host, cfg.Server.Port, &logger, &solveRecorder, &cfg)

	// start webserver
	logger.Info(fmt.Sprintf("Starting webserver on: %v:%v", server.URL, server.Port))
	err = http.ListenAndServe(
		fmt.Sprintf("%v:%v", server.URL, server.Port),
		server)
	if err != nil {
		logger.Fatal(nil, err, "Webserver stopped")
	}

}
