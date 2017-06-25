package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"markusreschke.name/selfhostedchatpolling/config"
	"markusreschke.name/selfhostedchatpolling/slack"
)

func parseSlashCommandRequest(appConfig config.AppConfig, logger *log.Logger, writer http.ResponseWriter, request *http.Request) (slack.SlashCommandRequest, error) {
	if appConfig.LogTraffic {
		logger.Printf("Poll Creation Request: %v\n", request)
	}
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return slack.SlashCommandRequest{}, errors.New("MethodNotAllowed")
	}
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return slack.SlashCommandRequest{}, errors.Wrap(err, "BadRequest - Body couldn't be read!")

	}
	if appConfig.LogTraffic {
		logger.Printf("Poll Creation Request Query String: %s\n", string(body))
	}
	parsedBody, err := url.ParseQuery(string(body))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return slack.SlashCommandRequest{}, errors.Wrap(err, "BadRequest - Body couldn't be parsed!")
	}
	slackRequest := slack.NewSlackRequest(parsedBody)

	if slackRequest.Token != appConfig.SlackVerificationToken {
		writer.WriteHeader(http.StatusUnauthorized)
		return slack.SlashCommandRequest{}, errors.New("Unauthorized")
	}

	return slackRequest, nil
}

func parseButtonActionRequest(appConfig config.AppConfig, logger *log.Logger, writer http.ResponseWriter, request *http.Request) (slack.ActionResponse, error) {
	if appConfig.LogTraffic {
		logger.Printf("Poll Update Request: %v\n", request)
	}
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return slack.ActionResponse{}, errors.New("MethodNotAllowed")
	}
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return slack.ActionResponse{}, errors.Wrap(err, "BadRequest - Couldn't read body")
	}
	parsedBody, err := url.ParseQuery(string(body))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return slack.ActionResponse{}, errors.Wrap(err, "BadRequest - Couldn't parse request")
	}
	payload := parsedBody.Get("payload")
	if payload == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return slack.ActionResponse{}, errors.New("BadRequest - No Payload")
	}
	if appConfig.LogTraffic {
		logger.Printf("Poll Update Request Payload: %s\n", payload)
	}

	actionCallback, err := slack.NewActionResponseFromPayload(payload)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return slack.ActionResponse{}, errors.Wrap(err, "InternalServerError - Error creating new action response from payload!")
	}

	if actionCallback.Token != appConfig.SlackVerificationToken {
		writer.WriteHeader(http.StatusUnauthorized)
		return slack.ActionResponse{}, errors.New("Unauthorized")
	}
	return actionCallback, nil
}
