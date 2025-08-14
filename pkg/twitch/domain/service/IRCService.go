package service

import (
	"fmt"
	"marluxGitHub/twitchbot/pkg/twitch/domain/model"
	"os"
	"os/signal"
	"syscall"

	"github.com/adeithe/go-twitch"
	"github.com/adeithe/go-twitch/irc"
)

type IRCService interface {
	Connect(string) error
}

type IRCServiceImpl struct {
	config *model.Config

	writer *irc.Conn
	reader *irc.Client
}

func (t *IRCServiceImpl) Connect(outhAuthToken string) error {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	writer := &irc.Conn{}
	writer.SetLogin(t.config.Twitch.Username, "oauth:"+outhAuthToken)

	if err := writer.Connect(); err != nil {
		panic("failed to start writer")
	}

	reader := twitch.IRC()

	if err := reader.Join(t.config.Twitch.Channel); err != nil {
		panic(err)
	}

	t.writer = writer
	t.reader = reader

	fmt.Println("Connected to IRC!")

	<-sc
	fmt.Println("Stopping...")
	reader.Close()
	writer.Close()

	return nil
}

func NewIRCService(config *model.Config) IRCService {
	return &IRCServiceImpl{config: config}
}
