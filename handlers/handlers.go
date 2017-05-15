package handlers

import (
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"markusreschke.name/selfhostedsimplepolling/config"
	"markusreschke.name/selfhostedsimplepolling/poll"
	"markusreschke.name/selfhostedsimplepolling/slack"
	"net/http"
	"net/url"
	"strings"
)

const (
	contentTypeJSON       = "application/json"
	httpHeaderContentType = "Content-Type"
)

func parseSlashCommand(commandArguments string) []string {
	return strings.Split(commandArguments, " ")
}

func GetNewPollRequestHandler(appConfig config.AppConfig, logger *log.Logger, pollStore poll.Store) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
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
		options := commandArguments[1:]
		question := commandArguments[0]
		response := slack.NewPollMessage(callBackID, question, options...)

		poll := poll.NewSimplePoll(callBackID.String(), commandArguments[0], slackRequest.UserID, options)

		pollStore.AddPoll(poll)

		writer.Header().Set(httpHeaderContentType, contentTypeJSON)
		writer.WriteHeader(http.StatusOK)

		responseJSON, _ := response.ToJSON()
		writer.Write(responseJSON)
		return
	}
}

func GetUpdatePollRequestHandler(appConfig config.AppConfig, logger *log.Logger, pollStore poll.Store) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
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

		results := pollStore.GetResult(actionCallback.CallbackID)

		poll := pollStore.GetPoll(actionCallback.CallbackID)

		updatedMessage := slack.UpdatePollMessage(poll, actionCallback, results)

		writer.Header().Set(httpHeaderContentType, contentTypeJSON)
		writer.WriteHeader(http.StatusOK)

		responseJSON, _ := updatedMessage.ToJSON()
		writer.Write(responseJSON)

		return
	}
}
