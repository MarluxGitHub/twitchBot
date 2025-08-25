package application

import (
	"marluxGitHub/twitchbot/pkg/twitch/domain/service"
	"os"
	"os/signal"
	"syscall"
)

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

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	t.apiService.SetAuthToken(oauthToken)
	t.apiService.Connect()

	go t.ircService.Connect(oauthToken)

	<-sc

	return nil
}

func NewChatBot() (ChatBot, error) {
	configService := service.NewConfigService()

	config, err := configService.LoadConfig()

	if err != nil {
		return nil, err
	}

	oauthService := service.NewOAuth2Service(
		config,
		configService)

	apiService := service.NewAPIService(config)
	ircService := service.NewIRCService(config, apiService)

	return &ChatBotImpl{
		configService: configService,
		oauthService:  oauthService,
		ircService:    ircService,
		apiService:    apiService,
	}, nil
}
