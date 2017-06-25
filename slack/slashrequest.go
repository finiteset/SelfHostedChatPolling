package slack

import (
	"net/url"
)

type SlashCommandRequest struct {
	Token       string
	TeamID      string
	TeamDomain  string
	ChannelID   string
	ChannelName string
	UserID      string
	UserName    string
	Command     string
	MsgText     string
	ResponseURL string
}

func NewSlackRequest(requestParams url.Values) SlashCommandRequest {
	request := SlashCommandRequest{}
	request.Token = requestParams.Get("token")
	request.TeamID = requestParams.Get("team_id")
	request.TeamDomain = requestParams.Get("team_domain")
	request.ChannelID = requestParams.Get("channel_id")
	request.ChannelName = requestParams.Get("channel_name")
	request.UserID = requestParams.Get("user_id")
	request.UserName = requestParams.Get("user_name")
	request.Command = requestParams.Get("command")
	request.MsgText = requestParams.Get("text")
	request.ResponseURL = requestParams.Get("response_url")
	return request
}
