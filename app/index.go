package main

import (
	"fmt"
	"github.com/dim13/gold"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

func (p Page) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pg := parsePage(*r.URL)
	app := conf.Blog.ArticlesPerPage
	p.Articles, p.NextPage, p.PrevPage = p.Articles.Page(pg, app)
	p.TagCloud = data.Articles.TagCloud()
	p.Config = conf

	err := tmpl.ExecuteTemplate(w, "index.tmpl", p)
	if err != nil {
		log.Fatal(err)
	}
}

func parsePage(u url.URL) int {
	if page, ok := u.Query()["page"]; ok {
		if pg, err := strconv.Atoi(page[0]); err == nil {
			return pg
		}
	}
	return 1
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var p Page
	var a gold.Articles

	switch {
	case strings.HasPrefix(r.URL.Path, "/tag/"):
		s := r.URL.Path[len("/tag/"):]
		a = data.Articles.Tag(s)
		p.Title = fmt.Sprint(conf.Blog.Title, " - ", s)
	case r.URL.Path == "/":
		a = data.Articles.Enabled()
		p.Title = conf.Blog.Title
	default:
		ar, err := data.Articles.Find(r.URL.Path[1:])
		if err == nil {
			p.Title = ar.Title
			a = gold.Articles{ar}
		}
	}

	if a == nil {
		http.NotFound(w, r)
		return
	}

	p.Articles = a
	p.ServeHTTP(w, r)
}

func assetHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

/*
func year(w http.ResponseWriter, r *http.Request) {
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

func month(w http.ResponseWriter, r *http.Request) {
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
*/
