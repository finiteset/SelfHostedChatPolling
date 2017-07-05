package slack

import (
	"bytes"
	"sort"
	"strings"

	"fmt"
	"strconv"

	slackApi "github.com/nlopes/slack"
	"markusreschke.name/selfhostedchatpolling/poll"
)

var RefreshButtonActionValue string = "refresh"
var PollDetailButtonActionValue string = "poll_details"
var ResponseTypeInChannel string = "in_channel"
var ResponseTypeEphemeral string = "ephemeral"

func NewVoteDetailMessage(results map[string][]string) SlackMessage {
	var messageText bytes.Buffer
	buildVoteDetailMessageTest(results, &messageText)
	slackMsg := SlackMessage{}
	slackMsg.ResponseType = ResponseTypeEphemeral
	slackMsg.Text = messageText.String()
	slackMsg.ReplaceOriginal = false
	return slackMsg
}

func buildVoteDetailMessageTest(results map[string][]string, messageText *bytes.Buffer) {
	var options []string
	for option := range results {
		options = append(options, option)
	}
	sort.Strings(options)
	for _, option := range options {
		messageText.WriteString("• ")
		messageText.WriteString(option)
		messageText.WriteString(": ")
		messageText.WriteString(strings.Join(results[option], ", "))
		messageText.WriteString("\n")
	}
}

func NewPollDetailButtonAttachment(poll poll.Poll) Attachment {
	var buttonAttachment Attachment
	buttonAttachment.Fallback = "Poll not available"
	buttonAttachment.CallbackID = poll.ID
	buttonAttachment.Color = "#0000ff"
	button := Action{PollDetailButtonActionValue + "_button", "Show vote details", "button", PollDetailButtonActionValue}
	buttonAttachment.AddAction(button)
	return buttonAttachment
}

func NewPollMessage(poll poll.Poll, results map[int]uint64) SlackMessage {
	var msg SlackMessage
	msg.ResponseType = ResponseTypeInChannel
	msg.Text = poll.Question
	msg.ReplaceOriginal = true
	var buttonAttachment Attachment
	for index, option := range poll.Options {
		var voteCount uint64 = results[index]
		buttonAttachment = Attachment{}
		buttonAttachment.Fallback = "Poll not available"
		buttonAttachment.CallbackID = poll.ID
		buttonAttachment.Text = option
		voteCountText := "Vote"
		if voteCount != 1 {
			voteCountText += "s"
		}
		var button Action
		button.Name = option + "_button"
		button.Text = fmt.Sprintf("%d %s", voteCount, voteCountText)
		button.Type = "button"
		button.Value = strconv.Itoa(index)
		buttonAttachment.AddAction(button)
		msg.AddAttachment(buttonAttachment)
	}
	if !poll.Anonymous {
		msg.AddAttachment(NewPollDetailButtonAttachment(poll))
	}
	msg.AddAttachment(NewRefreshButtonAttachment(poll))
	return msg
}

func NewRefreshButtonAttachment(poll poll.Poll) Attachment {
	var refreshButtonAttachment Attachment
	refreshButtonAttachment.Fallback = "Poll not available"
	refreshButtonAttachment.CallbackID = poll.ID
	refreshButton := Action{RefreshButtonActionValue + "_button", "Refresh", "button", RefreshButtonActionValue}
	refreshButtonAttachment.Color = "#ff0000"
	refreshButtonAttachment.AddAction(refreshButton)
	return refreshButtonAttachment
}

func ResolveVotersForPollDetails(pollDetails *map[string][]string, slackApiCient *slackApi.Client) error {
	for option, userList := range *pollDetails {
		resolvedUsers, err := ResolveVoterNamesBySlackID(&userList, slackApiCient)
		if err != nil {
			return err
		}
		(*pollDetails)[option] = resolvedUsers
	}
	return nil
}

func ResolveVoterNamesBySlackID(userIDs *[]string, slackApiClient *slackApi.Client) ([]string, error) {
	var resolvedUsers []string
	for _, userID := range *userIDs {
		resolvedUser, err := ResolveVoterNameBySlackID(userID, slackApiClient)
		if err != nil {
			return resolvedUsers, err
		}
		resolvedUsers = append(resolvedUsers, resolvedUser)
	}
	return resolvedUsers, nil
}

func ResolveVoterNameBySlackID(userID string, slackApiClient *slackApi.Client) (string, error) {
	user, err := slackApiClient.GetUserInfo(userID)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s (<@%s>)", user.RealName, userID), nil
}
