package service

import (
	"fmt"
	"marluxGitHub/twitchbot/pkg/twitch/model"
	"os"
	"os/signal"
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

	if err := reader.Join("MarluxMcLuxi"); err != nil {
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
	fmt.Printf("#%s %s: %s\n", msg.Channel, msg.Sender.DisplayName, msg.Text)

	if !msg.Sender.IsBroadcaster {
		t.writer.Say("#marluxmcluxi", msg.Sender.DisplayName+" mag DinoDance")
	}
}

func NewTwitchService(config *model.Config) TwitchService {
	return &TwitchServiceImpl{
		config: config,
	}
}
