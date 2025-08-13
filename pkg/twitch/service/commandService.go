package service

import (
	"strings"
	"time"

	"github.com/adeithe/go-twitch/irc"
)

var Commands = map[string]func(*TwitchServiceImpl, irc.ChatMessage){
	"hello": hello,
	"help":  help,
	"death": death,
	// Add more commands as needed
}

var CommandDescription = map[string]string{
	"!hello": "Says hello to the user.",
	"!help":  "Lists all available commands.",
	"!death": "Records a death of MarluxMcLuxi.",
	// Add more command descriptions as needed
}

var Deaths int = 0
var DeathCoolDownActive bool = false

func help(t *TwitchServiceImpl, msg irc.ChatMessage) {
	t.writer.Say(msg.Channel, "Available commands:")

	for cmd, desc := range CommandDescription {
		t.writer.Say(msg.Channel, strings.Join([]string{"- " + cmd, "-", desc}, " "))
	}
}

func hello(t *TwitchServiceImpl, msg irc.ChatMessage) {
	t.writer.Sayf(msg.Channel, "Hello there! Welcome %s to our Channel DinoDance", msg.Sender.DisplayName)
}

func death(t *TwitchServiceImpl, msg irc.ChatMessage) {
	if !DeathCoolDownActive {
		Deaths++
		t.writer.Sayf(msg.Channel, "Oh no MarluxMcLuxi died T.T. First Caller was %s. Total deaths: %d", msg.Sender.DisplayName, Deaths)

		DeathCoolDownActive = true

		go func() {
			time.Sleep(20 * time.Second)
			DeathCoolDownActive = false
		}()
	}
}
