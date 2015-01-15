package main

import (
	"log"
	"net/http"

	"github.com/dim13/gold/blog"
)

type rss struct {
	URL      string
	Title    string
	Subtitle string
	Articles blog.Articles
}

func rssHandler(w http.ResponseWriter, r *http.Request) {
	app := Conf.Blog.ArticlesPerPage
	a := Blog.Articles().Limit(app)

	rss := rss{
		URL:      "http://" + r.Host,
		Title:    Conf.Blog.Title,
		Subtitle: Conf.Blog.Subtitle,
		Articles: a,
	}
	err := tmpl.ExecuteTemplate(w, "rss.tmpl", rss)
	if err != nil {
		log.Println(err)
	}
}
