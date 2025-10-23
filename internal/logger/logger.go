package logger

import (
	"log"
)

// Logger singleton to be used across the application
type Logger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
}

var (
	global Logger = loggerObj{} 
)

func Log() Logger { return global }

// Create struct to implement the Logger interface and log logs with INFO or ERROR type.
type loggerObj struct{}
func (loggerObj) Infof(format string, args ...any)  { log.Printf("[INFO]: "+format, args...) }
func (loggerObj) Errorf(format string, args ...any) { log.Printf("[ERROR]: "+format, args...) }
