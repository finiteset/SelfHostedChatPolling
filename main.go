package main

import (
	"github.com/cloudfoundry-community/go-cfenv"
	"log"
	"markusreschke.name/selfhostedsimplepolling/config"
	"markusreschke.name/selfhostedsimplepolling/handlers"
	"markusreschke.name/selfhostedsimplepolling/poll"
	"net/http"
	"os"
	"strconv"
	"github.com/IBM-Bluemix/go-cloudant"
	"errors"
)

func getCloudantCredentialsFromEnv(cloudantServiceName string) (user, password string, err error) {
	appEnv, err := cfenv.Current()
	if err != nil {
		return "","",err
	}
	services := appEnv.Services
	cloudantService, err := services.WithName(cloudantServiceName)
	if err != nil {
		return "","",err
	}
	user, isUserSet := cloudantService.CredentialString("username")
	if !isUserSet {
		return "","", errors.New("Cloudant Username not found in CF Env!")
	}
	password, isPasswordSet := cloudantService.CredentialString("password")
	if !isPasswordSet {
		return "","", errors.New("Cloudant Password not found in CF Env!")
	}
	return user, password, nil
}

func main() {
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)
	appConfig, err := config.ReadConfigFromEnv()
	if err != nil {
		logger.Fatal("Error reading config file: ", err)
	}
	cloudantUser, cloudantPassword, err := getCloudantCredentialsFromEnv("shsp-cloudant")
	if err != nil {
		logger.Fatalf("Couldn't fetch CLoudant credentials: %v", err)
	}
	cloudantClient, err := cloudant.NewClient(cloudantUser, cloudantPassword)
	if err != nil {
		logger.Fatalf("Couldn't connect to Cloudant: %v", err)
	}
	pollStore, err := poll.NewCloudantStore(cloudantClient, os.Getenv("CLOUDANT_DB"))
	if err != nil {
		logger.Fatalf("Couldn't create poll store: %v", err)
	}
	logger.Println(pollStore)
	http.HandleFunc("/newpoll", handlers.GetNewPollRequestHandler(appConfig, logger, pollStore))
	http.HandleFunc("/updatepoll", handlers.GetUpdatePollRequestHandler(appConfig, logger, pollStore))
	http.ListenAndServe(":"+strconv.Itoa(appConfig.Port), nil)
}
