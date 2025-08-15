package service

import (
	"marluxGitHub/twitchbot/pkg/twitch/domain/model"
	"marluxGitHub/twitchbot/pkg/twitch/infrastructure/twitch/oauth2client"
)

type OAuth2Service interface {
	GetOAuth2Token() (string, error)
}

type OAuth2ServiceImpl struct {
	config        *model.Config
	configService ConfigService
	oauth2Client  oauth2client.OAuth2Client
}

func (t *OAuth2ServiceImpl) GetOAuth2Token() (string, error) {
	oAuthConf, err := t.configService.GetOAuth2Config()

	if err != nil {
		return "", err
	}

	if oAuthConf == nil {
		// Wenn die Konfigurationsdatei nicht existiert, erstelle sie
		oAuthConf = &model.OAuth{}

		oAuthResponse, err := t.oauth2Client.GetOAuth2Token()

		if err != nil {
			return "", err
		}

		t.configService.WriteOAuth2Config(oAuthResponse)

		return oAuthResponse.AccessToken, nil
	}

	oAuthConf, err = t.oauth2Client.RefreshOAuth2Token(oAuthConf)

	if err != nil {
		return "", err
	}

	t.configService.WriteOAuth2Config(oAuthConf)

	return oAuthConf.AccessToken, nil
}

func NewOAuth2Service(
	config *model.Config,
	configService ConfigService) OAuth2Service {
	return &OAuth2ServiceImpl{
		config:        config,
		configService: configService,
		oauth2Client:  oauth2client.NewOAuth2Client(config),
	}
}
