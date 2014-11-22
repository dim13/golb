package main

import (
	"log"
	"net/http"

	"github.com/dim13/gold/articles"
)

type rss struct {
	URL      string
	Title    string
	Subtitle string
	Articles articles.Articles
}

func rssHandler(w http.ResponseWriter, r *http.Request) {
	app := conf.Blog.ArticlesPerPage
	a := blog.Enabled().Articles().Limit(app)

	rss := rss{
		URL:      "http://" + r.Host,
		Title:    conf.Blog.Title,
		Subtitle: conf.Blog.Subtitle,
		Articles: a,
	}
	err := tmpl.ExecuteTemplate(w, "rss.tmpl", rss)
	if err != nil {
		log.Println(err)
	}
}
