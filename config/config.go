package config

import (
	"errors"
	"os"
	"strconv"
)

var Version string = "0.0.1"

type AppConfig struct {
	SlackVerificationToken string
	SlackOAuthToken        string
	Port                   int
	DbName                 string
	LogTraffic             bool
}

func ReadConfigFromEnv() (AppConfig, error) {
	var config AppConfig
	config.SlackVerificationToken = os.Getenv("SLACK_TOKEN")
	if config.SlackVerificationToken == "" {
		return config, errors.New("SLACK_TOKEN environment variable is not set!")
	}
	config.SlackOAuthToken = os.Getenv("SLACK_OAUTH_TOKEN")
	if config.SlackOAuthToken == "" {
		return config, errors.New("SLACK_OAUTH_TOKEN environment variable is not set!")
	}
	port, err := strconv.Atoi(os.Getenv("SHSP_PORT"))
	if err != nil {
		port, err = strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			port = 80
		}
	}
	config.Port = port
	config.DbName = os.Getenv("CLOUDANT_DB")
	config.LogTraffic, err = strconv.ParseBool(os.Getenv("SHSP_LOG_TRAFFIC"))
	if err != nil {
		config.LogTraffic = false
	}
	return config, nil
}
