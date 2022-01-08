package main

import (
	"csn/RESTendpoints"
	"csn/db"
	"log"
	"net/http"
)

type environment struct {
	welcome RESTendpoints.WelcomeModel
}

func main() {
	database, connectionError := db.InitDatabase()
	if connectionError != nil {
		return
	}

	env := &environment{welcome: RESTendpoints.WelcomeModel{
		DB: database,
	}}

	http.HandleFunc("/info", RESTendpoints.Info)
	http.HandleFunc("/", env.welcome.Welcome)
	err := http.ListenAndServe(":443", nil)
	if err != nil {
		log.Fatal("Unable to start web server!", err)
		return
	}
}
