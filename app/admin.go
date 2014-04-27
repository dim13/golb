// Test APP
package main

import (
	"log"
	"net/http"
)

func adminList(w http.ResponseWriter, r *http.Request, s []string) {
	p := Page{
		Config:   conf,
		Title:    "Admin interface",
		Articles: data.Articles,
	}
	err := tmpl.ExecuteTemplate(w, "admin.tmpl", p)
	if err != nil {
		log.Fatal(err)
	}
}

func adminSlug(w http.ResponseWriter, r *http.Request, s []string) {
	var p Page

	a, err := data.Articles.Find(s[0])
	if err != nil {
		p = Page{
			Config: conf,
			Error:  err,
		}
	} else {
		p = Page{
			Config:  conf,
			Title:   a.Title,
			Article: a,
		}
	}

	err = tmpl.ExecuteTemplate(w, "admin.tmpl", p)
	if err != nil {
		log.Fatal(err)
	}
}
