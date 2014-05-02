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

func NewRss() Rss {
	a := data.Articles.Enabled()
	app := conf.Blog.ArticlesPerPage

	return Rss{
		Url:      conf.Blog.Url,
		Title:    conf.Blog.Title,
		Subtitle: conf.Blog.Subtitle,
		Articles: a[:app],
	}
}

func (rss Rss) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "rss.tmpl", rss)
	if err != nil {
		log.Fatal(err)
	}
}
