package main

import (
	"csn/RESTendpoints"
	"csn/db"
	"log"
	"net/http"
)

type (
	Environment struct {
		Welcome RESTendpoints.WelcomeModel
	}
)

func main() {
	database, connectionError := db.InitDatabase()
	if connectionError != nil {
		return
	}

	env := &Environment{Welcome: RESTendpoints.WelcomeModel{
		DB: database,
	}}

	http.HandleFunc("/info", RESTendpoints.Info)
	http.HandleFunc("/", env.Welcome.Welcome)
	err := http.ListenAndServe(":443", nil)
	if err != nil {
		log.Fatal("Unable to start web server!", err)
		return
	}
}
