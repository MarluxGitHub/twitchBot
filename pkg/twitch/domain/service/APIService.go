package service

import (
	"marluxGitHub/twitchbot/pkg/twitch/domain/model"

	"github.com/adeithe/go-twitch/api"
)

type APIService interface {
	Connect() error
	SetAuthToken(token string)
}

type APIServiceImpl struct {
	config *model.Config

	oAuthToken string
	client     *api.Client
}

func (t *APIServiceImpl) SetAuthToken(token string) {
	t.oAuthToken = token
}

func (t *APIServiceImpl) Connect() error {
	t.client = api.New(t.config.Twitch.ClientID)

	return nil
}

func NewAPIService(config *model.Config) APIService {
	return &APIServiceImpl{config: config}
}
