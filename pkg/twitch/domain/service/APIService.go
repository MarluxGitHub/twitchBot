package service

import "marluxGitHub/twitchbot/pkg/twitch/domain/model"

type APIService interface {
	Connect(string) error
}

type APIServiceImpl struct {
	config *model.Config
}

func (t *APIServiceImpl) Connect(oAuthToken string) error {
	return nil
}

func NewAPIService(config *model.Config) APIService {
	return &APIServiceImpl{config: config}
}
