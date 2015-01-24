package main

import (
	"log"
	"net/http"

	"github.com/dim13/gold/blog"
	"github.com/dim13/gold/storage"
)

type rssPage struct {
	URL      string
	Blog     storage.Blog
	Articles blog.Articles
}

func (p rssPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app := Conf.Blog.ArticlesPerPage
	a := Blog.Articles().Limit(app)

	p = rssPage{
		URL:      "http://" + r.Host,
		Blog:     Conf.Blog,
		Articles: a,
	}
	err := tmpl.ExecuteTemplate(w, "rss.tmpl", p)
	if err != nil {
		log.Println(err)
	}
}
