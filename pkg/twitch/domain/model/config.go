package model

type Config struct {
	Twitch TwitchConfig `json:"twitch"`
}

type TwitchConfig struct {
	Username          string `json:"username"`
	ClientID          string `json:"clientId"`
	ClientSecret      string `json:"clientSecret"`
	AuthorizationCode string `json:"authorizationCode"`
	Channel           string `json:"channel"`
}
