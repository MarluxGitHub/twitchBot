package model

type Config struct {
	Twitch TwitchConfig `json:"twitch"`
}

type TwitchConfig struct {
	Username string `json:"username"`
	OAuth    string `json:"oauth"`
	Channel  string `json:"channel"`
}
