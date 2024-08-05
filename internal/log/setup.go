package log

import (
	"Ricochet-Robot-Solver/internal/config"
	"fmt"
	"log"
	"os"
	"time"
)

// Setup initializes the logger and Sentry
func Setup(cfg config.Config) (applicationLogger Logger, solverLogger SolveRecorder) {
	applicationLogger = setupApplicationLogger(cfg)
	solverLogger = setupSolverLogger(cfg)

	return applicationLogger, solverLogger
}

func setupApplicationLogger(cfg config.Config) (applicationLogger Logger) {
	// Log date and time, short file name, and line number
	logFlags := log.LstdFlags | log.Lshortfile

	logFilePath := fmt.Sprintf("%s/%s.log", cfg.Logging.Application.File.Path, getNewApplicationLogFileName(cfg))

	// Create log directory if not exists
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		err := os.MkdirAll(cfg.Logging.Application.File.Name, 0700)
		if err != nil {
			panic(err)
		}
	}
	// Open log file in append mode, create if not exists
	logFile, _ := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	// Create logger
	applicationLogger = Logger{
		infoLogger:  log.New(os.Stdout, "INFO:  ", logFlags),
		warnLogger:  log.New(os.Stdout, "WARN:  ", logFlags),
		errorLogger: log.New(os.Stdout, "ERROR: ", logFlags),
		fatalLogger: log.New(os.Stdout, "FATAL: ", logFlags),
		logFile:     logFile,
	}

	return applicationLogger
}

func setupSolverLogger(cfg config.Config) (solverRecorder SolveRecorder) {
	// Log date and time, short file name, and line number
	logFlags := log.LstdFlags | log.Lshortfile

	logFilePath := fmt.Sprintf("%s/%s.log", cfg.Logging.Solver.File.Path, getNewSolverLogFileName(cfg))

	// Create log directory if not exists
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		err := os.MkdirAll(cfg.Logging.Solver.File.Name, 0700)
		if err != nil {
			panic(err)
		}
	}
	// Open log file
	logFile, _ := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	// Create logger
	solverLogger := log.New(os.Stdout, "INFO:  ", logFlags)
	solverLogger.SetOutput(logFile)

	// Create SolveRecorder
	solverRecorder = SolveRecorder{
		recorder: solverLogger,
		logFile:  logFile,
	}

	return solverRecorder
}

func getNewApplicationLogFileName(cfg config.Config) string {
	if cfg.Logging.Rotation == "counter" {
		// check highest number in log directory "./logs"
		// return highest number + 1
		// if no files, return 1

		files, err := os.ReadDir(cfg.Logging.Application.File.Path)
		if err != nil {
			panic(err)
		}
		maxCounter := 0
		for _, file := range files {
			var counter int
			_, err := fmt.Sscanf(file.Name(), "%s-%d", cfg.Logging.Application.File.Name, &counter)
			if err == nil && counter > maxCounter {
				maxCounter = counter
			}
		}

		return fmt.Sprintf("%s-%d", cfg.Logging.Application.File.Name, maxCounter+1)
	} else {
		// check if log file exists with same date
		// if exists, return same name
		return fmt.Sprintf("%s-%s", cfg.Logging.Application.File.Name, time.Now().Format("2006-01-02"))
	}
}

func getNewSolverLogFileName(cfg config.Config) string {
	// check highest number in log directory "./logs/solver"
	// return highest number + 1
	// if no files, return 1
	files, err := os.ReadDir(cfg.Logging.Solver.File.Path)
	if err != nil {
		panic(err)
	}

	maxCounter := 0
	for _, file := range files {
		var counter int
		_, err := fmt.Sscanf(file.Name(), "%s-%d", cfg.Logging.Solver.File.Name, &counter)
		if err == nil && counter > maxCounter {
			maxCounter = counter
		}
	}

	return fmt.Sprintf("%s-%d", cfg.Logging.Solver.File.Name, maxCounter+1)
}

func (r *SolveRecorder) UpdateLogFile(cfg config.Config) {
	// close current log file
	err := r.logFile.Close()
	if err != nil {
		panic(err)
	}

	// get new log file name
	logFileName := getNewSolverLogFileName(cfg)

	// update log file
	logFilePath := fmt.Sprintf("%s/%s.log", cfg.Logging.Solver.File.Path, logFileName)
	logFile, _ := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	r.recorder.SetOutput(logFile)
	r.logFile = logFile
}
