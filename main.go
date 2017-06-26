package main

import (
	"errors"
	"log"

	"github.com/cloudfoundry-community/go-cfenv"
	"markusreschke.name/selfhostedchatpolling/config"
	"markusreschke.name/selfhostedchatpolling/handlers"
	"markusreschke.name/selfhostedchatpolling/poll"
	//"markusreschke.name/selfhostedchatpolling/poll/memstore"
	"net/http"
	"os"
	"strconv"

	"github.com/IBM-Bluemix/go-cloudant"
	"markusreschke.name/selfhostedchatpolling/poll/cloudantstore"
)

func getCloudantCredentialsFromEnv(cloudantServiceName string) (user, password string, err error) {
	appEnv, err := cfenv.Current()
	if err != nil {
		return "", "", err
	}
	services := appEnv.Services
	cloudantService, err := services.WithName(cloudantServiceName)
	if err != nil {
		return "", "", err
	}
	user, isUserSet := cloudantService.CredentialString("username")
	if !isUserSet {
		return "", "", errors.New("cloudant username not found in CF env")
	}
	password, isPasswordSet := cloudantService.CredentialString("password")
	if !isPasswordSet {
		return "", "", errors.New("cloudant password not found in CF env")
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
		logger.Fatalf("Couldn't fetch Cloudant credentials: %v", err)
	}
	cloudantClient, err := cloudant.NewClient(cloudantUser, cloudantPassword)
	if err != nil {
		logger.Fatalf("Couldn't connect to Cloudant: %v", err)
	}
	pollStoreBackend, err := cloudantstore.NewCloudantStoreBackend(cloudantClient, appConfig.DbName)
	if err != nil {
		logger.Fatalf("Couldn't create poll store: %v", err)
	}
	//pollStoreBackend := memstore.NewInMemoryStoreBackend()
	pollStore := poll.NewDefaultStore(pollStoreBackend)
	http.HandleFunc("/newpoll", handlers.GetNewPollRequestHandler(appConfig, logger, pollStore, false))
	http.HandleFunc("/newpollanon", handlers.GetNewPollRequestHandler(appConfig, logger, pollStore, true))
	http.HandleFunc("/updatepoll", handlers.GetPollButtonRequestHandler(appConfig, logger, pollStore))
	http.HandleFunc("/version", handlers.GetVersionRequestHandler(appConfig, logger))
	http.ListenAndServe(":"+strconv.Itoa(appConfig.Port), nil)
}
