package oauth2client

import (
	"encoding/json"
	"io"
	"marluxGitHub/twitchbot/pkg/twitch/domain/model"
	"net/http"
	"strings"
)

type OAuth2Client interface {
	GetOAuth2Token() (string, error)
}

type OAuth2ClientImpl struct {
	config *model.Config
}

func (t *OAuth2ClientImpl) GetOAuth2Token() (string, error) {
	url := "https://id.twitch.tv/oauth2/token"

	oauth2Request := OAuth2Request{
		ClientID:     t.config.Twitch.ClientID,
		ClientSecret: t.config.Twitch.ClientSecret,
		Code:         t.config.Twitch.AutorisationCode,
		GrantType:    "authorization_code",
		RedirectURI:  "http://localhost",
	}

	// Die Daten für die Anfrage als x-www-form-urlencoded kodieren
	data := make(map[string]string)
	data["client_id"] = oauth2Request.ClientID
	data["client_secret"] = oauth2Request.ClientSecret
	data["code"] = oauth2Request.Code
	data["grant_type"] = oauth2Request.GrantType
	data["redirect_uri"] = oauth2Request.RedirectURI

	form := ""
	for k, v := range data {
		if form != "" {
			form += "&"
		}
		form += k + "=" + v
	}

	req, err := http.NewRequest("POST", url, io.NopCloser(strings.NewReader(form)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var oauth2Response OAuth2Response
	err = json.Unmarshal(body, &oauth2Response)

	if err != nil {
		return "", err
	}

	return oauth2Response.AccessToken, nil
}

func NewOAuth2Client(config *model.Config) OAuth2Client {
	return &OAuth2ClientImpl{config: config}
}
