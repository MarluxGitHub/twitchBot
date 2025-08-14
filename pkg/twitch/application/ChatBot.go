package application

import "marluxGitHub/twitchbot/pkg/twitch/domain/service"

type ChatBot interface {
	Start() error
}

type ChatBotImpl struct {
	ircService    service.IRCService
	apiService    service.APIService
	configService service.ConfigService
	oauthService  service.OAuth2Service
}

func (t *ChatBotImpl) Start() error {

	oauthToken, err := t.oauthService.GetOAuth2Token()

	if err != nil {
		return err
	}

	go t.ircService.Connect(oauthToken)
	go t.apiService.Connect(oauthToken)

	return nil
}

func NewChatBot() (ChatBot, error) {
	configService := service.NewConfigService()

	config, err := configService.LoadConfig()

	if err != nil {
		return nil, err
	}

	oauthService := service.NewOAuth2Service(config)
	ircService := service.NewIRCService(config)
	apiService := service.NewAPIService(config)

	return &ChatBotImpl{
		configService: configService,
		oauthService:  oauthService,
		ircService:    ircService,
		apiService:    apiService,
	}, nil
}
