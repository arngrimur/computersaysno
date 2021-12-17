package main

import (
	"RESTendpoints"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/info", RESTendpoints.Info)
	http.HandleFunc("/", RESTendpoints.Welcome)
	err := http.ListenAndServe(":443", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
