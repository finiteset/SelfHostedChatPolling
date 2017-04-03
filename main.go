package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	configFilePath        = "./appconfig.json"
	contentTypeJSON       = "application/json"
	httpHeaderContentType = "Content-Type"
)

type AppConfig struct {
	SlackVerificationToken string
}

type SlackRequest struct {
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

type SlackMessage struct {
	Text string `json:"text"`
}

func (m *SlackMessage) ToJson() ([]byte, error) {
	return json.Marshal(m)
}

var config AppConfig
var logger *log.Logger

func main() {
	logger = log.New(os.Stdout, "logger: ", log.Lshortfile)
	config = readConfig()
	http.HandleFunc("/", handleRequests)
	http.ListenAndServe(":8080", nil)
}

func readConfig() AppConfig {
	var parsedConfig AppConfig
	rawConfig, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Println("Couldn't load config:")
		os.Exit(1)
	}
	err = json.Unmarshal(rawConfig, &parsedConfig)
	return parsedConfig
}

func handleRequests(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		logger.Println("MethodNotAllowed")
		return
	}
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		logger.Println("BadRequest")
		return
	}
	parsedBody, err := url.ParseQuery(string(body))
	slackRequest := NewSlackRequest(parsedBody)
	logger.Println(fmt.Sprintf("%+v", slackRequest))

	if slackRequest.Token != config.SlackVerificationToken {
		writer.WriteHeader(http.StatusUnauthorized)
		logger.Println("Unauthorized")
		return
	} else {
		writer.WriteHeader(http.StatusOK)
		response := SlackMessage{}
		response.Text = fmt.Sprintf("%+v", slackRequest)
		responseJSON, _ := response.ToJson()
		logger.Println(fmt.Sprintf("JSON: %s", string(responseJSON)))
		writer.Write(responseJSON)
		return
	}
}

func NewSlackRequest(requestParams url.Values) SlackRequest {
	request := SlackRequest{}
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
