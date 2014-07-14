package main

import (
	"log"
	"net/http"
)

type Robots struct {
	Host string
}

func robotsHandler(w http.ResponseWriter, r *http.Request) {
	rob := Robots{
		Host: r.Host,
	}
	err := tmpl.ExecuteTemplate(w, "robots.tmpl", rob)
	if err != nil {
		log.Println(err)
	}
}
