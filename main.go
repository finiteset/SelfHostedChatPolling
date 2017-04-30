package main

import (
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"markusreschke.name/selfhostedsimplepolling/config"
	"markusreschke.name/selfhostedsimplepolling/poll"
	"markusreschke.name/selfhostedsimplepolling/slack"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	contentTypeJSON       = "application/json"
	httpHeaderContentType = "Content-Type"
)

var appConfig config.AppConfig
var logger *log.Logger
var pollStore poll.Store

func newPollMessage(callbackID uuid.UUID, question string, options ...string) slack.SlackMessage {
	var msg slack.SlackMessage
	msg.Text = question
	var buttonAttachment slack.Attachment
	buttonAttachment.Fallback = "Poll not available"
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

func updatePollMessage(poll poll.Poll, callback slack.ActionResponse, results map[string]uint64) slack.SlackMessage {
	var msg slack.SlackMessage
	msg.Text = poll.Question()
	var buttonAttachment slack.Attachment
	buttonAttachment.Fallback = "Poll not available"
	buttonAttachment.CallbackID = poll.ID()
	for index, option := range poll.Options() {
		var button slack.Action
		button.Name = option + "_button"
		button.Text = option + " " + fmt.Sprintf("%d", results[strconv.Itoa(index)])
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
	appConfig, err = config.ReadConfigFromEnv()
	if err != nil {
		logger.Fatal("Error reading config file: ", err)
	}
	pollStore = poll.NewInMemoryStore()
	logger.Println(pollStore)
	http.HandleFunc("/newpoll", handleNewPollRequests)
	http.HandleFunc("/updatepoll", handleUpdatePollRequests)
	http.ListenAndServe(":"+strconv.Itoa(appConfig.Port), nil)
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

	callBackID := uuid.NewV4()
	commandArguments := parseSlashCommand(slackRequest.MsgText)
	response := newPollMessage(callBackID, commandArguments[0], commandArguments[1:]...)

	poll := poll.NewSimplePoll(callBackID.String(), commandArguments[0], slackRequest.UserID)

	pollStore.AddPoll(poll)

	writer.Header().Set(httpHeaderContentType, contentTypeJSON)
	writer.WriteHeader(http.StatusOK)

	responseJSON, _ := response.ToJSON()
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

	vote := poll.NewSimpleVote(uuid.NewV4().String(), actionCallback.User.ID, actionCallback.CallbackID, actionCallback.Actions[0].Value)

	pollStore.AddVote(vote)

	logger.Println(pollStore)
	logger.Println(pollStore.GetResult(actionCallback.CallbackID))

	writer.Header().Set(httpHeaderContentType, contentTypeJSON)
	writer.WriteHeader(http.StatusOK)

	return
}
