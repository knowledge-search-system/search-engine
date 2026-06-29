package config

import (
	"os"

	"go.uber.org/fx"
)

var Module = fx.Module("config",
	fx.Provide(New),
)

func New() (*Config, error) {
	path := os.Getenv("SEARCH_ENGINE_CONFIG_PATH")
	if path == "" {
		path = "config/configs/config.yaml"
	}
	return Load(path)
}
