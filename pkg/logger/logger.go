package logger

import (
	"go.uber.org/zap"
)

var log = zap.NewNop()

func Init(cfg Config) error {
	var err error

	if cfg.Development {
		log, err = zap.NewDevelopment()
	} else {
		log, err = zap.NewProduction()
	}

	return err
}

func Sync() error {
	return log.Sync()
}

func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

func Fatal(msg string, err error) {
	log.Fatal(msg, zap.Error(err))
}
