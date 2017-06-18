package handlers

import (
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"markusreschke.name/selfhostedsimplepolling/config"
	"markusreschke.name/selfhostedsimplepolling/poll"
	"markusreschke.name/selfhostedsimplepolling/slack"
	"net/http"
	"net/url"
	"strconv"
)

const (
	contentTypeJSON       = "application/json"
	contentTypeText       = "text/plain"
	httpHeaderContentType = "Content-Type"
)

func parseNewPollRequest(appConfig config.AppConfig, logger *log.Logger, writer http.ResponseWriter, request *http.Request) (slack.SlackRequest, error) {
	if appConfig.LogTraffic {
		logger.Printf("Poll Creation Request: %v\n", request)
	}
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return slack.SlackRequest{}, errors.New("MethodNotAllowed")
	}
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return slack.SlackRequest{}, errors.Wrap(err, "BadRequest - Body couldn't be read!")

	}
	if appConfig.LogTraffic {
		logger.Printf("Poll Creation Request Query String: %s\n", string(body))
	}
	parsedBody, err := url.ParseQuery(string(body))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return slack.SlackRequest{}, errors.Wrap(err, "BadRequest - Body couldn't be parsed!")
	}
	slackRequest := slack.NewSlackRequest(parsedBody)

	if slackRequest.Token != appConfig.SlackVerificationToken {
		writer.WriteHeader(http.StatusUnauthorized)
		return slack.SlackRequest{}, errors.New("Unauthorized")
	}

	return slackRequest, nil
}

func handleUserFacingError(logger *log.Logger, writer http.ResponseWriter, err error, logMessage, slackMessage string) {
	logger.Println(logMessage, err)
	writer.Header().Set(httpHeaderContentType, contentTypeJSON)
	writer.WriteHeader(http.StatusOK) //Can't set another status than OK because slack then shows an error msg. of its own
	errorMsg := slack.NewSlackErrorMessage(slackMessage)
	errorResponseJSON, _ := errorMsg.ToJSON()
	writer.Write(errorResponseJSON)
}

func GetNewPollRequestHandler(appConfig config.AppConfig, logger *log.Logger, pollStore poll.Store) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		slackRequest, err := parseNewPollRequest(appConfig, logger, writer, request)
		if err != nil {
			logger.Println("Error reading and parsing request for new poll: ", err)
			return
		}
		callBackID := uuid.NewV4()
		commandArguments := slack.ParseSlashCommand(slackRequest.MsgText)
		options := commandArguments[1:]
		question := commandArguments[0]

		poll := poll.Poll{callBackID.String(), question, slackRequest.UserID, options}
		err = pollStore.AddPoll(poll)

		if err != nil {
			handleUserFacingError(logger, writer, err, "Error adding poll to store: ", "Error creating new poll!")
			return
		}

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
		actionCallback, err := parseUpdatePollRequest(appConfig, logger, writer, request)
		if err != nil {
			logger.Println("Error reading and parsing request for poll update: ", err)
			return
		}

		actionValue := actionCallback.Actions[0].Value

		if actionValue != slack.RefreshButtonActionValue {
			voteOptionIndex, err := strconv.Atoi(actionValue)
			if err != nil {
				writer.WriteHeader(http.StatusBadRequest)
				logger.Println("BadRequest - Value of Action Callback is not a valid vote option index", err)
				return
			}
			vote := poll.Vote{uuid.NewV4().String(), actionCallback.User.ID, actionCallback.CallbackID, voteOptionIndex}
			err = pollStore.AddVote(vote)
			if err != nil {
				handleUserFacingError(logger, writer, err, "Error adding vote to store: ", "Error submitting vote!")
				return
			}
		}

		results, err := pollStore.GetResult(actionCallback.CallbackID)
		if err != nil {
			handleUserFacingError(logger, writer, err, "Error calculating current poll count: ", "Error refreshing poll!")
			return
		}

		poll, err := pollStore.GetPoll(actionCallback.CallbackID)
		if err != nil {
			handleUserFacingError(logger, writer, err, "Error fetching poll from store for message recreation: ", "Error refreshing poll!")
			return
		}

		updatedMessage := slack.NewPollMessage(poll, results)

		writer.Header().Set(httpHeaderContentType, contentTypeJSON)
		writer.WriteHeader(http.StatusOK)

		responseJSON, _ := updatedMessage.ToJSON()
		writer.Write(responseJSON)

		return
	}
}

func parseUpdatePollRequest(appConfig config.AppConfig, logger *log.Logger, writer http.ResponseWriter, request *http.Request) (slack.ActionResponse, error) {
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
