package main

import (
	"github.com/dim13/gold"
	"log"
	"net/http"
)

type AdminPage struct {
	Articles gold.Articles
	Article  *gold.Article
	Title    string
	Config   *gold.Config
	Error    string
}

func (p AdminPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.Config = conf
	err := tmpl.ExecuteTemplate(w, "admin.tmpl", p)
	if err != nil {
		log.Fatal(err)
	}
}

type AdminIndex struct { AdminPage }

func (p *AdminIndex) Select(match []string) {
	p.Articles = data.Articles
	p.Title = "Admin Interface"
}

type AdminSlug struct { AdminPage }

func (p *AdminSlug) Select(match []string) {
	a, err := data.Articles.Find(match[0])
	if err == nil {
		p.Title = a.Title
		p.Article = a
	}
}

/*
import (
	"log"
	"net/http"
)

func adminListHandler(w http.ResponseWriter, r *http.Request, s []string) {
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

func adminSlugHandler(w http.ResponseWriter, r *http.Request, s []string) {
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
*/
