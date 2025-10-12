package utils

import (
	"log"
	"os"
)

type Logger struct {
	INFO  *log.Logger
	WARN  *log.Logger
	ERROR *log.Logger
}

var LOGGER Logger

func(Logger) Init() {
	LOGGER = Logger{
		INFO:  log.New(os.Stdout, "[INFO]: ", log.Ldate|log.Ltime|log.Lshortfile),
		WARN:  log.New(os.Stdout, "[WARN]: ", log.Ldate|log.Ltime|log.Lshortfile),
		ERROR: log.New(os.Stdout, "[ERROR]: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
