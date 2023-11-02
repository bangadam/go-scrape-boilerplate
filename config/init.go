package config

import (
	"github.com/spf13/viper"
)

func Init() error {
	// config path in config/config.json
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	return viper.ReadInConfig()
}
