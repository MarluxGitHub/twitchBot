package main

import (
	"marluxGitHub/twitchbot/pkg/twitch/service"
)

func main() {
	// Initialize ConfigService
	configService := service.NewConfigLoader()

	config, err := configService.LoadConfig()

	if err != nil {
		panic(err)
	}

	// Initialize the Twitch service
	twitchService := service.NewTwitchService(config)
	// Connect to Twitch
	if err := twitchService.Connect(); err != nil {
		panic(err)
	}
}
