package main

import "marluxGitHub/twitchbot/pkg/twitch/application"

func main() {
	// Initialize ConfigService
	chatBot, err := application.NewChatBot()

	if err != nil {
		panic(err)
	}

	chatBot.Start()
}
