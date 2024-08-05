package log

import (
	"os"
	// External packages
)

// Info logs an info message to the console and log file
func (l *Logger) Info(message ...interface{}) {
	errorAndMessage := ErrorAndMessage{Err: nil, Message: message}

	// Write to console
	_ = l.infoLogger.Output(2, errorAndMessage.Error())

	// Write to log file
	l.infoLogger.SetOutput(l.logFile)
	_ = l.infoLogger.Output(2, errorAndMessage.Error())
	l.infoLogger.SetOutput(os.Stdout)
}

// Warn logs a warning to the console and log file, and captures it with Sentry
func (l *Logger) Warn(err error, message ...interface{}) {
	errorAndMessage := ErrorAndMessage{Err: err, Message: message}

	// Write to console
	_ = l.warnLogger.Output(2, errorAndMessage.Error())

	// Write to log file
	l.warnLogger.SetOutput(l.logFile)
	_ = l.warnLogger.Output(2, errorAndMessage.Error())
	l.warnLogger.SetOutput(os.Stdout)
}

// Error logs an error to the console and log file, and captures it with Sentry
func (l *Logger) Error(err error, message ...interface{}) {
	errorAndMessage := ErrorAndMessage{Err: err, Message: message}

	// Write to console
	_ = l.errorLogger.Output(2, errorAndMessage.Error())

	// Write to log file
	l.errorLogger.SetOutput(l.logFile)
	_ = l.errorLogger.Output(2, errorAndMessage.Error())
	l.errorLogger.SetOutput(os.Stdout)
}

// Fatal logs a fatal error to the console and log file, and captures it with Sentry
func (l *Logger) Fatal(err error, message ...interface{}) {
	errorAndMessage := ErrorAndMessage{Err: err, Message: message}

	// Write to console
	_ = l.fatalLogger.Output(2, errorAndMessage.Error())

	// Write to log file
	l.fatalLogger.SetOutput(l.logFile)
	_ = l.fatalLogger.Output(2, errorAndMessage.Error())
	l.fatalLogger.SetOutput(os.Stdout)
}

// Record solver output
func (r *SolveRecorder) Record(message ...interface{}) {
	errorAndMessage := ErrorAndMessage{Err: nil, Message: message}

	_ = r.recorder.Output(2, errorAndMessage.Error())
}
