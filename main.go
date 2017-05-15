package main

import (
	"log"
	"markusreschke.name/selfhostedsimplepolling/config"
	"markusreschke.name/selfhostedsimplepolling/handlers"
	"markusreschke.name/selfhostedsimplepolling/poll"
	"net/http"
	"os"
	"strconv"
)

func main() {
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)
	appConfig, err := config.ReadConfigFromEnv()
	if err != nil {
		logger.Fatal("Error reading config file: ", err)
	}
	pollStore := poll.NewInMemoryStore()
	logger.Println(pollStore)
	http.HandleFunc("/newpoll", handlers.GetNewPollRequestHandler(appConfig, logger, pollStore))
	http.HandleFunc("/updatepoll", handlers.GetUpdatePollRequestHandler(appConfig, logger, pollStore))
	http.ListenAndServe(":"+strconv.Itoa(appConfig.Port), nil)
}
