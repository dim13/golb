package main

import (
	"fmt"
	"github.com/dim13/gold"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Page struct {
	Config   *gold.Config
	Title    string
	Articles gold.Articles
	Article  *gold.Article
	Error    error
	PrevPage int
	NextPage int
	TagCloud gold.TagCloud
}

func index(w http.ResponseWriter, r *http.Request, s []string) {
	var p Page

	a, err := data.Articles.Find(s[0])
	if err == nil {
		p = Page{
			Title:    a.Title,
			Articles: gold.Articles{a},
		}
		genpage(w, p)
	} else {
		page(w, r, []string{"1"})
	}
}

func page(w http.ResponseWriter, r *http.Request, s []string) {
	pg, err := strconv.Atoi(s[0])
	if err != nil {
		log.Fatal(err)
	}
	app := conf.Blog.ArticlesPerPage

	a, next, prev := data.Articles.Enabled().Page(pg, app)

	p := Page{
		Title:    conf.Blog.Title,
		Articles: a,
		NextPage: next,
		PrevPage: prev,
	}

	genpage(w, p)
}

func tags(w http.ResponseWriter, r *http.Request, s []string) {
	p := Page{
		Title:    fmt.Sprint(conf.Blog.Title, " - ", s[0]),
		Articles: data.Articles.Tag(s[0]),
	}
	genpage(w, p)
}

func assets(w http.ResponseWriter, r *http.Request, s []string) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func year(w http.ResponseWriter, r *http.Request, s []string) {
	y, err := strconv.Atoi(s[0])
	if err != nil {
		log.Fatal(err)
	}
	a := data.Articles.Year(y)
	p := Page{
		Title:    fmt.Sprint(conf.Blog.Title, " - ", y),
		Articles: a,
	}
	genpage(w, p)
}

func month(w http.ResponseWriter, r *http.Request, s []string) {
	y, err := strconv.Atoi(s[0])
	if err != nil {
		log.Fatal(err)
	}
	m, err := strconv.Atoi(s[1])
	if err != nil {
		log.Fatal(err)
	}

	a := data.Articles.Enabled().Year(y).Month(m)
	p := Page{
		Title:    fmt.Sprint(conf.Blog.Title, " - ", y, time.Month(m)),
		Articles: a,
	}
	genpage(w, p)
}

func genpage(w http.ResponseWriter, p Page) {
	p.TagCloud = data.Articles.TagCloud()
	p.Config = conf

	err := tmpl.ExecuteTemplate(w, "index.tmpl", p)
	if err != nil {
		log.Fatal(err)
	}
}
