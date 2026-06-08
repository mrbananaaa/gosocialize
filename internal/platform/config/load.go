package config

import (
	"github.com/joho/godotenv"
	"github.com/mrbananaaa/gosocialize/pkg/cfgloader"
)

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg, err := cfgloader.LoadCfg[Config]()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
