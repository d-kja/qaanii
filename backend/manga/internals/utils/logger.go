package utils

import (
	"log"

	"go.uber.org/zap"
)

type Logger struct {
	Instance *zap.Logger
	Sugar    *zap.SugaredLogger
}

var logger Logger = Logger{}

func (Logger) SetupLogger() Logger {
	instance, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Unable to create logger instance, error: %+v\n", err)
	}

	sugar := instance.Sugar()
	logger = Logger{
		Instance: instance,
		Sugar:    sugar,
	}

	return logger
}

func GetLogger() *zap.SugaredLogger {
	return logger.Sugar
}
