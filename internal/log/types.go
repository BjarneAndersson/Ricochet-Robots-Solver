package log

import (
	"fmt"
	// Standard library packages
	"log"
	"os"
	// External packages
	// Internal packages
)

// Logger inherits all basic log operations
type Logger struct {
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
	logFile     *os.File
}

type SolveRecorder struct {
	recorder *log.Logger
	logFile  *os.File
}

// ErrorAndMessage Create a custom error for Sentry
type ErrorAndMessage struct {
	Err     error
	Message []interface{}
}

// Error Transform the error and message into output
func (e *ErrorAndMessage) Error() string {
	// Convert interface to string
	msg := ""
	for _, v := range e.Message {
		msg += fmt.Sprintf("%v", v)
	}

	// Format output based on error existence
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", msg, e.Err)
	}
	return fmt.Sprintf("%s", msg)
}
