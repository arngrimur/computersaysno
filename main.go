package main

import (
	"csn/RESTendpoints"
	"csn/db"
	"log"
	"net/http"
)

type (
	Environment struct {
		welcome RESTendpoints.WelcomeModel
	}
)

func main() {
	database, connectionError := db.InitDatabase()
	if connectionError != nil {
		return
	}

	env := &Environment{welcome: RESTendpoints.WelcomeModel{
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
