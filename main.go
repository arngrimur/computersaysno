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
	db, dbErr := db.Init("testuser:testpassword@(localhost:3306)/csn_db")
	if dbErr != nil {
		return
	}
	defer db.Close()

	env := &environment{welcome: RESTendpoints.WelcomeModel{
		DB: db,
	}}

	http.HandleFunc("/info", RESTendpoints.Info)
	http.HandleFunc("/", env.welcome.Welcome)
	err := http.ListenAndServe(":443", nil)
	if err != nil {
		log.Fatal("Unable to start web server!", err)
		return
	}
}
