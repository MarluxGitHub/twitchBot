package service

import (
	"marluxGitHub/twitchbot/pkg/twitch/domain/model"

	"github.com/spf13/viper"
)

type ConfigService interface {
	LoadConfig() (*model.Config, error)
}

type ConfigServiceImpl struct {
}

func (t *ConfigServiceImpl) LoadConfig() (*model.Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../../config")

	var config model.Config
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func NewConfigService() ConfigService {
	return &ConfigServiceImpl{}
}
