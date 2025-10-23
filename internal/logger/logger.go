package logger

import (
	"log"
)

type Logger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
}

var (
	global Logger = loggerObj{} 
)

func Log() Logger { return global }


type loggerObj struct{}
func (loggerObj) Infof(format string, args ...any)  { log.Printf("[INFO]: "+format, args...) }
func (loggerObj) Errorf(format string, args ...any) { log.Printf("[ERROR]: "+format, args...) }
