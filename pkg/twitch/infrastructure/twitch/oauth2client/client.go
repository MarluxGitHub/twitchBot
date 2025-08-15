package oauth2client

import (
	"encoding/json"
	"io"
	"marluxGitHub/twitchbot/pkg/twitch/domain/model"
	"net/http"
	"strings"
)

const (
	oauth2TokenUrl = "https://id.twitch.tv/oauth2/token"
)

type OAuth2Client interface {
	GetOAuth2Token() (*model.OAuth, error)
	RefreshOAuth2Token(*model.OAuth) (*model.OAuth, error)
}

type OAuth2ClientImpl struct {
	config *model.Config
}

func (t *OAuth2ClientImpl) formEncode(data map[string]string) string {
	form := ""
	for k, v := range data {
		if form != "" {
			form += "&"
		}
		form += k + "=" + v
	}

	return form
}

func (t *OAuth2ClientImpl) RefreshOAuth2Token(oauth *model.OAuth) (*model.OAuth, error) {
	// Die Daten für die Anfrage als x-www-form-urlencoded kodieren
	data := make(map[string]string)
	data["client_id"] = t.config.Twitch.ClientID
	data["client_secret"] = t.config.Twitch.ClientSecret
	data["grant_type"] = "refresh_token"
	data["refresh_token"] = oauth.RefreshToken

	form := t.formEncode(data)

	return doRequest(form)

}

func (t *OAuth2ClientImpl) GetOAuth2Token() (*model.OAuth, error) {
	// Die Daten für die Anfrage als x-www-form-urlencoded kodieren
	data := make(map[string]string)
	data["client_id"] = t.config.Twitch.ClientID
	data["client_secret"] = t.config.Twitch.ClientSecret
	data["code"] = t.config.Twitch.AutorisationCode
	data["grant_type"] = "authorization_code"
	data["redirect_uri"] = "http://localhost"

	form := t.formEncode(data)

	return doRequest(form)

}

func doRequest(form string) (*model.OAuth, error) {
	req, err := http.NewRequest("POST", oauth2TokenUrl, io.NopCloser(strings.NewReader(form)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var oauth2Response OAuth2Response
	err = json.Unmarshal(body, &oauth2Response)

	if err != nil {
		return nil, err
	}

	oauth := model.OAuth{
		AccessToken:  oauth2Response.AccessToken,
		RefreshToken: oauth2Response.RefreshToken,
		ExpiresIn:    oauth2Response.ExpiresIn,
	}

	return &oauth, nil
}

func NewOAuth2Client(config *model.Config) OAuth2Client {
	return &OAuth2ClientImpl{config: config}
}
