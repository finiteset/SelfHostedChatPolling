package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"markusreschke.name/selfhostedsimplepolling/slack"
	"markusreschke.name/selfhostedsimplepolling/config"
)

const (
	configFilePath        = "./appconfig.json"
	contentTypeJSON       = "application/json"
	httpHeaderContentType = "Content-Type"
)

var appConfig config.AppConfig
var logger *log.Logger

func main() {
	logger = log.New(os.Stdout, "logger: ", log.Lshortfile)
	var err error;
	appConfig, err = config.ReadConfig(configFilePath)
	if (err != nil) {
		logger.Fatal("Error reading config file: ", err)
	}
	http.HandleFunc("/", handleRequests)
	http.ListenAndServe(":8080", nil)
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
	slackRequest := slack.NewSlackRequest(parsedBody)
	logger.Println(fmt.Sprintf("%+v", slackRequest))

	if slackRequest.Token != appConfig.SlackVerificationToken {
		writer.WriteHeader(http.StatusUnauthorized)
		logger.Println("Unauthorized")
		return
	}
	writer.Header().Set(httpHeaderContentType, contentTypeJSON)
	writer.WriteHeader(http.StatusOK)
	response := slack.SlackMessage{}
	response.Text = fmt.Sprintf("%+v", slackRequest)
	responseJSON, _ := response.ToJSON()
	logger.Println(fmt.Sprintf("JSON: %s", string(responseJSON)))
	writer.Write(responseJSON)
	return
}



