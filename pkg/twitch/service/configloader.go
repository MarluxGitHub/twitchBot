package service

import (
	"marluxGitHub/twitchbot/pkg/twitch/model"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type ConfigLoader interface {
	LoadConfig() (*model.Config, error)
}

type ConfigLoaderImpl struct {
}

func (c *ConfigLoaderImpl) LoadConfig() (*model.Config, error) {
	wd, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	wd = strings.TrimSuffix(wd, "\\cmd\\bot")

	configPath := filepath.Join(wd, "config")
	viper.AddConfigPath(configPath) // path to look for the config file in (Projektroot als Start)
	viper.SetConfigName("config")   // name of config file (without extension)
	viper.SetConfigType("json")     // REQUIRED if the config file does not have the extension in the name

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config model.Config

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func NewConfigLoader() ConfigLoader {
	return &ConfigLoaderImpl{}
}
