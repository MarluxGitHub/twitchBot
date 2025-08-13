package service

import (
	"fmt"
	"marluxGitHub/twitchbot/pkg/twitch/model"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/adeithe/go-twitch"
	"github.com/adeithe/go-twitch/irc"
)

type TwitchService interface {
	Connect() error
}

type TwitchServiceImpl struct {
	config *model.Config
	writer *irc.Conn
}

func (t *TwitchServiceImpl) Connect() error {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	writer := &irc.Conn{}
	writer.SetLogin(t.config.Twitch.Username, "oauth:"+t.config.Twitch.OAuth)

	t.writer = writer

	if err := writer.Connect(); err != nil {
		panic("failed to start writer")
	}

	reader := twitch.IRC()

	reader.OnShardReconnect(onShardReconnect)
	reader.OnShardLatencyUpdate(onShardLatencyUpdate)
	reader.OnShardMessage(t.onShardMessage)

	if err := reader.Join(t.config.Twitch.Channel); err != nil {
		panic(err)
	}

	fmt.Println("Connected to IRC!")

	<-sc
	fmt.Println("Stopping...")
	reader.Close()
	writer.Close()

	return nil
}

func onShardReconnect(shardID int) {
	fmt.Printf("Shard #%d reconnected\n", shardID)
}

func onShardLatencyUpdate(shardID int, latency time.Duration) {
	fmt.Printf("Shard #%d has %dms ping\n", shardID, latency.Milliseconds())
}

func (t *TwitchServiceImpl) onShardMessage(shardID int, msg irc.ChatMessage) {
	// Handle incoming chat messages with !
	if strings.HasPrefix(msg.Text, "!") {
		// regex first word after !
		re := regexp.MustCompile(`^!(\w+)`)
		matches := re.FindStringSubmatch(msg.Text)

		if len(matches) > 1 {
			command := matches[1]

			// Check if Commands Dictionary contains a entry for command
			if cmd, ok := Commands[command]; ok {
				cmd(t, msg)
			} else {
				t.writer.Say(msg.Channel, fmt.Sprintf("Unknown command: %s", command))
			}
		}
	}
}

func NewTwitchService(config *model.Config) TwitchService {
	return &TwitchServiceImpl{
		config: config,
	}
}
