package config

import "github.com/spf13/viper"

type Config struct {
	HOST string `mapstruct:"host"`
	DSN  string `mapstruct:"dsn"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./pkg/config/envs")
	viper.SetConfigName("dev")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
