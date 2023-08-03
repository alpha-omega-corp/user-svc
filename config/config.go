package config

import (
	"embed"
	"github.com/alpha-omega-corp/services/config"
	"io/fs"
	"sync"
)

var (
	//go:embed envs
	embedFS      embed.FS
	unwrapFSOnce sync.Once
	unwrappedFS  fs.FS
)

func Config(env string) (*config.Config, error) {
	return config.LoadConfig(config.MakeConfig(
		&unwrapFSOnce,
		embedFS,
		unwrappedFS,
	), env)
}
