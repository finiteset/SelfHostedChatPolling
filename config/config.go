package config

import (
	"errors"
	"os"
	"strconv"
)

type AppConfig struct {
	SlackVerificationToken string
	Port                   int
}

func ReadConfigFromEnv() (AppConfig, error) {
	var config AppConfig
	config.SlackVerificationToken = os.Getenv("SLACK_TOKEN")
	if config.SlackVerificationToken == "" {
		return config, errors.New("SLACK_TOKEN environment variable is not set!")
	}
	port, err := strconv.Atoi(os.Getenv("SHSP_PORT"))
	if err != nil {
		port, err = strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			port = 80
		}
	}
	config.Port = port
	return config, nil
}
