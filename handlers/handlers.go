package handlers

import (
	"log"
	"net/http"
	"strconv"

	slackApi "github.com/nlopes/slack"
	"github.com/satori/go.uuid"
	"markusreschke.name/selfhostedchatpolling/config"
	"markusreschke.name/selfhostedchatpolling/poll"
	"markusreschke.name/selfhostedchatpolling/slack"
)

const (
	contentTypeJSON       = "application/json"
	contentTypeText       = "text/plain"
	httpHeaderContentType = "Content-Type"
)

func handleUserFacingError(logger *log.Logger, writer http.ResponseWriter, err error, logMessage, slackMessage string) {
	logger.Println(logMessage, err)
	writer.Header().Set(httpHeaderContentType, contentTypeJSON)
	writer.WriteHeader(http.StatusOK) //Can't set another status than OK because slack then shows an error msg. of its own
	errorMsg := slack.NewSlackErrorMessage(slackMessage)
	errorResponseJSON, _ := errorMsg.ToJSON()
	writer.Write(errorResponseJSON)
}

func GetNewPollRequestHandler(appConfig config.AppConfig, logger *log.Logger, pollStore poll.Store, forAnonPolls bool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		slackRequest, err := parseSlashCommandRequest(appConfig, logger, writer, request)
		if err != nil {
			logger.Println("Error reading and parsing request for new poll: ", err)
			return
		}
		callBackID := uuid.NewV4()
		commandArguments := slack.ParseSlashCommand(slackRequest.MsgText)
		options := commandArguments[1:]
		question := commandArguments[0]

		poll := poll.Poll{ID: callBackID.String(), Question: question, CreatorID: slackRequest.UserID, Options: options, Anonymous: forAnonPolls}
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

func GetPollButtonRequestHandler(appConfig config.AppConfig, logger *log.Logger, pollStore poll.Store) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		actionCallback, err := parseButtonActionRequest(appConfig, logger, writer, request)
		if err != nil {
			logger.Println("Error reading and parsing request for poll update: ", err)
			return
		}

		actionValue := actionCallback.Actions[0].Value

		switch actionValue {
		case slack.PollDetailButtonActionValue:
			writePollDetailMessage(writer, logger, pollStore, actionCallback, appConfig)
		case slack.RefreshButtonActionValue:
			writeUpdatedPollMessage(writer, logger, pollStore, actionCallback)
		default:
			handleNewVoteRequest(writer, logger, pollStore, actionCallback)
			writeUpdatedPollMessage(writer, logger, pollStore, actionCallback)
		}
	}
}

func handleNewVoteRequest(writer http.ResponseWriter, logger *log.Logger, pollStore poll.Store, actionCallback slack.ActionResponse) {
	logger.Println("Handle new vote request")
	voteOptionIndex, err := strconv.Atoi(actionCallback.Actions[0].Value)
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

func writeUpdatedPollMessage(writer http.ResponseWriter, logger *log.Logger, pollStore poll.Store, actionCallback slack.ActionResponse) {
	logger.Println("Handle poll update")
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
}

func writePollDetailMessage(writer http.ResponseWriter, logger *log.Logger, pollStore poll.Store, actionCallback slack.ActionResponse, appConfig config.AppConfig) {
	logger.Println("Handle poll detail request")
	voteDetails, err := pollStore.GetVoteDetails(actionCallback.CallbackID)
	if err != nil {
		handleUserFacingError(logger, writer, err, "Error calculating current poll vote details: ", "Error getting poll details!")
		return
	}

	err = slack.ResolveVotersForPollDetails(&voteDetails, slackApi.New(appConfig.SlackOAuthToken))
	if err != nil {
		handleUserFacingError(logger, writer, err, "Error getting user names from Slack API: ", "Error getting user names for voters!")
		return
	}

	pollDetailMessage := slack.NewVoteDetailMessage(voteDetails)

	writer.Header().Set(httpHeaderContentType, contentTypeJSON)
	writer.WriteHeader(http.StatusOK)

	responseJSON, _ := pollDetailMessage.ToJSON()
	writer.Write(responseJSON)
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
