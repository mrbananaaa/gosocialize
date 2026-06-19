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
