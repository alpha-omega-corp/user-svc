package config

import (
	"github.com/alpha-omega-corp/services/config"
	"github.com/alpha-omega-corp/services/server"
	"github.com/spf13/viper"
)

var (
	v        = viper.New()
	cManager = server.NewConfigManager(v)
)

func HostsConfig() (config.HostsConfig, error) {
	return cManager.HostsConfig()
}

func AuthConfig() (config.AuthenticationConfig, error) {
	return cManager.AuthConfig()
}
