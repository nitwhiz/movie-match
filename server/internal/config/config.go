package config

import "github.com/spf13/viper"

func Init() error {
	viper.SetEnvPrefix("MM")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	return viper.ReadInConfig()
}
