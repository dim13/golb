package main

import (
	"log"
	"net/http"

	"github.com/dim13/gold/blog"
)

type rssPage struct {
	URL         string
	Title       string
	Description string
	Articles    blog.Articles
}

func (p rssPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app := Conf.Blog.ArticlesPerPage
	a := Blog.Articles().Limit(app)

	p = rssPage{
		URL:         "http://" + r.Host,
		Title:       Conf.Blog.Title,
		Description: Conf.Blog.Description,
		Articles:    a,
	}
	err := tmpl.ExecuteTemplate(w, "rss.tmpl", p)
	if err != nil {
		log.Println(err)
	}
}
