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
	dbConfig := db.DbConfig{
		DbSecrets: db.DbSecrets{
			RootPassword: "secret",
			MysqlUser:    "testuser",
			MysqlPwd:     "testpassword",
		},
		HostConfig: db.HostConfig{
			AutoRemove:    false,
			RestartPolicy: "always",
		},
	}
	connString, pool, resource := db.SetupDatbase(dbConfig)
	defer db.Purge(pool, resource)
	database, connectionError := db.InitDatabase(*connString)
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
