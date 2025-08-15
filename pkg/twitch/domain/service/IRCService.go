package service

import (
	"fmt"
	"marluxGitHub/twitchbot/pkg/twitch/domain/model"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/adeithe/go-twitch"
	"github.com/adeithe/go-twitch/irc"
)

type IRCService interface {
	Connect(string) error
}

type IRCServiceImpl struct {
	config         *model.Config
	commandService CommandService

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

	reader.OnShardChannelJoin(t.onShardChannelJoin)
	reader.OnShardMessage(t.onShardRawMessage)

	<-sc
	fmt.Println("Stopping...")
	reader.Close()
	writer.Close()

	return nil
}

func (t *IRCServiceImpl) onShardChannelJoin(shard int, channel, user string) {
	t.writer.Sayf(t.config.Twitch.Channel, "Welcome @%s DinoDance", user)

}

func (t *IRCServiceImpl) onShardRawMessage(shard int, message irc.ChatMessage) {
	if message.Sender.IsModerator {
		return
	}

	if strings.HasPrefix(message.Text, "!") {
		t.commandService.HandleCommand(message.Text)
	}
}

func NewIRCService(
	config *model.Config,
) IRCService {
	return &IRCServiceImpl{
		config:         config,
		commandService: NewCommandService(),
	}
}
