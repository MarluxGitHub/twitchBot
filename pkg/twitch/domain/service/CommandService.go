package service

import (
	"strings"

	"github.com/adeithe/go-twitch/irc"
)

var deathCounter int = 0
var isDeathOnCooldown bool = false

var commands map[string]func(IRCService, irc.ChatMessage) = map[string]func(IRCService, irc.ChatMessage){
	"lurk": func(ircService IRCService, message irc.ChatMessage) {
		ircService.SendMessage("#"+message.Channel, "User is now lurking.")
	},
}

type CommandService interface {
	HandleCommand(ircService IRCService, message irc.ChatMessage)
}

type CommandServiceImpl struct {
}

func (t *CommandServiceImpl) HandleCommand(ircService IRCService, message irc.ChatMessage) {
	text := message.Text
	text = strings.TrimPrefix(text, "!") // Remove the '!' prefix

	// get from text
	command := strings.Split(text, " ")[0]

	if cmdFunc, ok := commands[command]; ok {
		cmdFunc(ircService, message)
	} else {
		ircService.SendMessage("#"+message.Channel, "Unknown command: "+command)
	}
}

func NewCommandService() CommandService {
	return &CommandServiceImpl{}
}
