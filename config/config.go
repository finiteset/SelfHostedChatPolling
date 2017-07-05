package config

import (
	"errors"
	"os"
	"strconv"
)

const (
	BackendCloudant = "cloudant"
	BackendInMemory = "inmemory"
)

var (
	Version string = "v1.0.0"
)

type AppConfig struct {
	SlackVerificationToken string
	SlackOAuthToken        string
	Port                   int
	DbName                 string
	LogTraffic             bool
	Backend                string
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
	config.Backend = os.Getenv("SHCP_BACKEND")
	if config.Backend == "" {
		config.Backend = BackendInMemory
	}
	port, err := strconv.Atoi(os.Getenv("SHCP_PORT"))
	if err != nil {
		port, err = strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			port = 80
		}
	}
	config.Port = port
	config.DbName = os.Getenv("CLOUDANT_DB")
	config.LogTraffic, err = strconv.ParseBool(os.Getenv("SHCP_LOG_TRAFFIC"))
	if err != nil {
		config.LogTraffic = false
	}
	return config, nil
}
