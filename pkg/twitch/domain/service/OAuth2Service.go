package service

import (
	"marluxGitHub/twitchbot/pkg/twitch/domain/model"
	"marluxGitHub/twitchbot/pkg/twitch/infrastructure/twitch/oauth2client"
)

type OAuth2Service interface {
	GetOAuth2Token() (string, error)
}

type OAuth2ServiceImpl struct {
	config       *model.Config
	oauth2Client oauth2client.OAuth2Client
}

func (t *OAuth2ServiceImpl) GetOAuth2Token() (string, error) {
	token, err := t.oauth2Client.GetOAuth2Token()

	if err != nil {
		return "", err
	}

	return token, nil
}

func NewOAuth2Service(config *model.Config) OAuth2Service {
	return &OAuth2ServiceImpl{
		config:       config,
		oauth2Client: oauth2client.NewOAuth2Client(config),
	}
}
