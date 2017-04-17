package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"markusreschke.name/selfhostedsimplepolling/config"
	"markusreschke.name/selfhostedsimplepolling/slack"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	configFilePath        = "./appconfig.json"
	contentTypeJSON       = "application/json"
	httpHeaderContentType = "Content-Type"
)

var appConfig config.AppConfig
var logger *log.Logger

func newPollMessage(question string, options ...string) slack.SlackMessage {
	var msg slack.SlackMessage
	msg.Text = question
	var buttonAttachment slack.Attachment
	buttonAttachment.Fallback = "Poll not available"
	callbackID := uuid.NewV4()
	buttonAttachment.CallbackID = callbackID.String()
	for index, option := range options {
		var button slack.Action
		button.Name = option + "_button"
		button.Text = option
		button.Type = "button"
		button.Value = strconv.Itoa(index)
		buttonAttachment.AddAction(button)
	}
	msg.AddAttachment(buttonAttachment)
	return msg
}

func main() {
	logger = log.New(os.Stdout, "logger: ", log.Lshortfile)
	var err error
	appConfig, err = config.ReadConfig(configFilePath)
	if err != nil {
		logger.Fatal("Error reading config file: ", err)
	}
	http.HandleFunc("/newpoll", handleNewPollRequests)
	http.HandleFunc("/updatepoll", handleUpdatePollRequests)
	http.ListenAndServe(":8080", nil)
}

func parseSlashCommand(commandArguments string) []string {
	return strings.Split(commandArguments, " ")
}

func handleNewPollRequests(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		logger.Println("MethodNotAllowed")
		return
	}
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		logger.Println("BadRequest", err)
		return
	}
	parsedBody, err := url.ParseQuery(string(body))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		logger.Println("BadRequest", err)
		return
	}
	slackRequest := slack.NewSlackRequest(parsedBody)

	if slackRequest.Token != appConfig.SlackVerificationToken {
		writer.WriteHeader(http.StatusUnauthorized)
		logger.Println("Unauthorized")
		return
	}

	commandArguments := parseSlashCommand(slackRequest.MsgText)
	response := newPollMessage(commandArguments[0], commandArguments[1:]...)

	writer.Header().Set(httpHeaderContentType, contentTypeJSON)
	writer.WriteHeader(http.StatusOK)

	responseJSON, _ := response.ToJSON()
	logger.Println(fmt.Sprintf("JSON: %s", string(responseJSON)))
	writer.Write(responseJSON)
	return
}

func handleUpdatePollRequests(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		logger.Println("MethodNotAllowed")
		return
	}
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		logger.Println("BadRequest", err)
		return
	}
	parsedBody, err := url.ParseQuery(string(body))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		logger.Println("BadRequest", err)
		return
	}
	payload := parsedBody.Get("payload")
	if payload == "" {
		writer.WriteHeader(http.StatusBadRequest)
		logger.Println("BadRequest")
		return
	}

	actionCallback, err := slack.NewActionResponseFromPayload(payload)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		logger.Println("InternalServerError", err)
		return
	}

	if actionCallback.Token != appConfig.SlackVerificationToken {
		writer.WriteHeader(http.StatusUnauthorized)
		logger.Println("Unauthorized")
		return
	}

	writer.Header().Set(httpHeaderContentType, contentTypeJSON)
	writer.WriteHeader(http.StatusOK)
	return
}
