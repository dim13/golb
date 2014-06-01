package main

import (
	"github.com/dim13/gold/articles"
	"log"
	"net/http"
)

type Rss struct {
	Url      string
	Title    string
	Subtitle string
	Articles articles.Articles
}

func rssHandler(w http.ResponseWriter, r *http.Request) {
	a := data.Articles.Enabled()
	app := conf.Blog.ArticlesPerPage
	rss := Rss{
		Url:      "http://" + r.Host,
		Title:    conf.Blog.Title,
		Subtitle: conf.Blog.Subtitle,
		Articles: a[:app],
	}
	err := tmpl.ExecuteTemplate(w, "rss.tmpl", rss)
	if err != nil {
		log.Fatal(err)
	}
}
