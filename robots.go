package main

import (
	"log"
	"net/http"
)

type robots struct {
	URL string
}

func robotsHandler(w http.ResponseWriter, r *http.Request) {
	rob := robots{
		URL: r.Host,
	}
	err := tmpl.ExecuteTemplate(w, "robots.tmpl", rob)
	if err != nil {
		log.Println(err)
	}
}
