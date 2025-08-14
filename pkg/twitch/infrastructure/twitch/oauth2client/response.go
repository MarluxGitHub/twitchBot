package oauth2client

// OAuth2Response repräsentiert die Antwortstruktur des Twitch OAuth2-Tokens
type OAuth2Response struct {
	AccessToken  string   `json:"access_token"`
	ExpiresIn    int      `json:"expires_in"`
	RefreshToken string   `json:"refresh_token"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"token_type"`
}
