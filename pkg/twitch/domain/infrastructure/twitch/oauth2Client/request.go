package oauth2client

// OAuth2Request repräsentiert die Struktur für eine OAuth2-Token-Anfrage an Twitch
type OAuth2Request struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	GrantType    string `json:"grant_type"`
	RedirectURI  string `json:"redirect_uri"`
}
