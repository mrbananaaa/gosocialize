package logger

import (
	"sync"

	"go.uber.org/zap"
)

var once sync.Once
var log = zap.NewNop()

func Init(cfg Config) error {
	var err error

	once.Do(func() {
		if cfg.Development {
			logConfig := zap.NewDevelopmentConfig()
			logConfig.DisableStacktrace = true

			log, err = logConfig.Build()
		} else {
			log, err = zap.NewProduction()
		}
	})

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
