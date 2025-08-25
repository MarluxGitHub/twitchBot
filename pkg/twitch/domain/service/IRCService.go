package service

import (
	"fmt"
	"marluxGitHub/twitchbot/pkg/twitch/domain/model"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/adeithe/go-twitch"
	"github.com/adeithe/go-twitch/irc"
)

type IRCService interface {
	Connect(string) error
	SendMessage(channel, message string)
}

type IRCServiceImpl struct {
	config         *model.Config
	commandService CommandService
	apiService     APIService

	writer      *irc.Conn
	reader      *irc.Client
	stopWelcome chan struct{}
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

	reader.OnShardMessage(t.onShardRawMessage)

	t.stopWelcome = make(chan struct{})
	go t.botWelcomeMessage()

	<-sc
	close(t.stopWelcome)
	fmt.Println("Stopping...")
	reader.Close()
	writer.Close()

	return nil
}

func (t *IRCServiceImpl) onShardRawMessage(shard int, message irc.ChatMessage) {
	if message.Sender.IsModerator {
		return
	}

	if strings.HasPrefix(message.Text, "!") {
		t.commandService.HandleCommand(t, message)

	}
}

func (t *IRCServiceImpl) botWelcomeMessage() {
	for {
		select {
		case <-t.stopWelcome:
			return
		default:
			t.writer.Sayf(t.config.Twitch.Channel, "Welcome on the Channel. I'm ModMilli on the Duty. If you like what you see, feel free to follow! DinoDance")
			// 10 Minuten warten oder früher abbrechen
			select {
			case <-t.stopWelcome:
				return
			case <-time.After(10 * time.Minute):
			}
		}
	}

}

func (t *IRCServiceImpl) SendMessage(channel, message string) {
	t.writer.Say(channel, message)
}

func NewIRCService(
	config *model.Config,
	apiService APIService,
) IRCService {
	return &IRCServiceImpl{
		config:         config,
		commandService: NewCommandService(),
		apiService:     apiService,
	}
}
