package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type TMDBProviderConfig struct {
	APIKey        string `mapstructure:"api_key"`
	PosterBaseUrl string `mapstructure:"poster_base_url"`
	Language      string `mapstructure:"language"`
	Region        string `mapstructure:"region"`
}

type MediaProvidersConfig struct {
	TMDB *TMDBProviderConfig `mapstructure:"tmdb,omitempty"`
}

type PosterConfig struct {
	FsBasePath string `mapstructure:"fs_base_path"`
}

type UserConfig struct {
	Username    string `mapstructure:"username"`
	DisplayName string `mapstructure:"display_name"`
	Password    string `mapstructure:"password"`
}

type LoginConfig struct {
	JWTKey string       `mapstructure:"jwt_key"`
	Users  []UserConfig `mapstructure:"users"`
}

type Config struct {
	Database       DatabaseConfig       `mapstructure:"database"`
	MediaProviders MediaProvidersConfig `mapstructure:"media_providers"`
	PosterConfig   PosterConfig         `mapstructure:"poster"`
	Login          LoginConfig          `mapstructure:"login"`
	Debug          bool                 `mapstructure:"debug"`
}

var C Config

func Init() error {
	viper.SetEnvPrefix("MOVIEMATCH")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()

	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warn("config file not found")
		} else {
			return err
		}
	}

	if err := viper.Unmarshal(&C); err != nil {
		return err
	}

	return nil
}
