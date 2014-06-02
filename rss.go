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
	app := conf.Blog.ArticlesPerPage
	a := data.Articles.Enabled().Limit(app)

	rss := Rss{
		Url:      "http://" + r.Host,
		Title:    conf.Blog.Title,
		Subtitle: conf.Blog.Subtitle,
		Articles: a,
	}
	err := tmpl.ExecuteTemplate(w, "rss.tmpl", rss)
	if err != nil {
		log.Fatal(err)
	}
}
