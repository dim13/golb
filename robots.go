package main

import (
	"log"
	"net/http"
)

type robotsPage struct {
	URL string
}

func (p robotsPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.URL = "http://" + r.Host
	err := tmpl.ExecuteTemplate(w, "robots.tmpl", p)
	if err != nil {
		log.Println(err)
	}
}
