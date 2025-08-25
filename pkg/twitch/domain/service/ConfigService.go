package service

import (
	"encoding/json"
	"fmt"
	"marluxGitHub/twitchbot/pkg/twitch/domain/model"
	"os"

	"github.com/spf13/viper"
)

const (
	oAuthConfFile = "oAuthConf.json"
)

var filepath string = oAuthConfFile

type ConfigService interface {
	LoadConfig() (*model.Config, error)
	GetOAuth2Config() (*model.OAuth, error)
	WriteOAuth2Config(oauth2 *model.OAuth) error
}

type ConfigServiceImpl struct {
}

func (t *ConfigServiceImpl) LoadConfig() (*model.Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../../config")

	var config model.Config
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (t *ConfigServiceImpl) GetOAuth2Config() (*model.OAuth, error) {
	viper.SetConfigFile(filepath)

	if err := viper.ReadInConfig(); err != nil {
		// Wenn die Konfigurationsdatei nicht existiert, teste anderen Pfad
		filepath = "cmd/bot/" + oAuthConfFile
		viper.SetConfigFile(filepath)

		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("config not found")

			// Wenn die Konfigurationsdatei nicht existiert, gib nil zurück
			return nil, nil
		}
	}

	var oauth2 model.OAuth

	if err := viper.Unmarshal(&oauth2); err != nil {
		return nil, err
	}

	return &oauth2, nil
}

func (t *ConfigServiceImpl) WriteOAuth2Config(oauth2 *model.OAuth) error {
	// Schreibe die OAuth2-Daten im gewünschten Format in die oAuthConf.json
	viper.SetConfigType("json")
	viper.SetConfigFile(filepath)

	// Erstelle eine Map, die dem gewünschten JSON-Format entspricht
	data := map[string]any{
		"accessToken":  oauth2.AccessToken,
		"refreshToken": oauth2.RefreshToken,
		"expiresIn":    oauth2.ExpiresIn,
	}

	// Schreibe die Map direkt in die Datei
	// Wir verwenden hier die Standardbibliothek, um das Format exakt zu steuern
	file, err := os.Create(oAuthConfFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

func NewConfigService() ConfigService {
	return &ConfigServiceImpl{}
}
