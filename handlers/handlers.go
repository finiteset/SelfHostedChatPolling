package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/satori/go.uuid"
	"markusreschke.name/selfhostedsimplepolling/config"
	"markusreschke.name/selfhostedsimplepolling/poll"
	"markusreschke.name/selfhostedsimplepolling/slack"
)

const (
	contentTypeJSON       = "application/json"
	contentTypeText       = "text/plain"
	httpHeaderContentType = "Content-Type"
)

func GetNewPollRequestHandler(appConfig config.AppConfig, logger *log.Logger, pollStore poll.Store) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if appConfig.LogTraffic {
			logger.Printf("Poll Creation Request: %v\n", request)
		}
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
		if appConfig.LogTraffic {
			logger.Printf("Poll Creation Request Query String: %s\n", string(body))
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
		commandArguments := slack.ParseSlashCommand(slackRequest.MsgText)
		options := commandArguments[1:]
		question := commandArguments[0]

		poll := poll.Poll{callBackID.String(), question, slackRequest.UserID, options}
		pollStore.AddPoll(poll)

		response := slack.NewPollMessage(poll, nil)

		writer.Header().Set(httpHeaderContentType, contentTypeJSON)
		writer.WriteHeader(http.StatusOK)

		responseJSON, _ := response.ToJSON()
		writer.Write(responseJSON)
		return
	}
}

func GetUpdatePollRequestHandler(appConfig config.AppConfig, logger *log.Logger, pollStore poll.Store) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if appConfig.LogTraffic {
			logger.Printf("Poll Update Request: %v\n", request)
		}
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
			logger.Println("BadRequest - Couldn't parse request", err)
			return
		}
		payload := parsedBody.Get("payload")
		if payload == "" {
			writer.WriteHeader(http.StatusBadRequest)
			logger.Println("BadRequest - No Payload")
			return
		}
		if appConfig.LogTraffic {
			logger.Printf("Poll Update Request Payload: %s\n", payload)
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

		voteOptionIndex, err := strconv.Atoi(actionCallback.Actions[0].Value)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			logger.Println("BadRequest - Value of Action Callback is not a valid vote option index", err)
			return
		}
		vote := poll.Vote{uuid.NewV4().String(), actionCallback.User.ID, actionCallback.CallbackID, voteOptionIndex}

		pollStore.AddVote(vote)

		results, _ := pollStore.GetResult(actionCallback.CallbackID)

		poll, _ := pollStore.GetPoll(actionCallback.CallbackID)

		updatedMessage := slack.NewPollMessage(poll, results)

		writer.Header().Set(httpHeaderContentType, contentTypeJSON)
		writer.WriteHeader(http.StatusOK)

		responseJSON, _ := updatedMessage.ToJSON()
		writer.Write(responseJSON)

		return
	}
}
func GetVersionRequestHandler(appConfig config.AppConfig, logger *log.Logger) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if appConfig.LogTraffic {
			logger.Printf("Version Request: %v\n", request)
		}
		if request.Method != http.MethodGet {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			logger.Println("MethodNotAllowed")
			return
		}

		writer.Header().Set(httpHeaderContentType, contentTypeText)
		writer.WriteHeader(http.StatusOK)

		writer.Write([]byte(config.Version))

		return
	}
}
