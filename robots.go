package main

import (
	"log"
	"net/http"
)

type Robots struct {
	Url string
}

func robotsHandler(w http.ResponseWriter, r *http.Request) {
	rob := Robots{
		Url: "http://" + r.Host,
	}
	err := tmpl.ExecuteTemplate(w, "robots.tmpl", rob)
	if err != nil {
		log.Fatal(err)
	}
}
