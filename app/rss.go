package main

import (
	"github.com/dim13/gold"
	"log"
	"net/http"
)

type Rss struct {
	Url      string
	Title    string
	Subtitle string
	Articles gold.Articles
}

func rss(w http.ResponseWriter, r *http.Request, s []string) {
	a := data.Articles.Enabled()
	app := conf.Blog.ArticlesPerPage

	p := Rss{
		Url:      conf.Blog.Url,
		Title:    conf.Blog.Title,
		Subtitle: conf.Blog.Subtitle,
		Articles: a[:app],
	}

	err := tmpl.ExecuteTemplate(w, "rss.tmpl", p)
	if err != nil {
		log.Fatal(err)
	}
}
