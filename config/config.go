package config

import (
	"encoding/json"
	"io/ioutil"
)

type AppConfig struct {
	SlackVerificationToken string
}

func ReadConfig(configFilePath string) (AppConfig, error) {
	var parsedConfig AppConfig
	rawConfig, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return parsedConfig, err
	}
	err = json.Unmarshal(rawConfig, &parsedConfig)
	return parsedConfig, err
}