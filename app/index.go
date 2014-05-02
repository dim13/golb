package main

import (
	"fmt"
	"github.com/dim13/gold"
	"log"
	"net/http"
	"net/url"
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
	Match    []string
	Year     int
	Month    time.Month
}

func atoiMust(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func parsePage(u url.URL) int {
	if page, ok := u.Query()["page"]; ok {
		if pg, err := strconv.Atoi(page[0]); err == nil {
			return pg
		}
	}
	return 1
}

func (p *Page) StoreMatch(s []string) { p.Match = s }

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

type TagPage struct{ Page }

func (p TagPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s := p.Page.Match[0]
	p.Articles = data.Articles.Tag(s)
	p.Title = fmt.Sprint(conf.Blog.Title, " - ", s)
	p.Page.ServeHTTP(w, r)
}

type IndexPage struct{ Page }

func (p IndexPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.Articles = data.Articles.Enabled()
	p.Title = conf.Blog.Title
	p.Page.ServeHTTP(w, r)
}

type SlugPage struct{ Page }

func (p SlugPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a, err := data.Articles.Find(p.Match[0])
	if err != nil {
		http.NotFound(w, r)
		return
	}
	p.Title = a.Title
	p.Articles = gold.Articles{a}
	p.Page.ServeHTTP(w, r)
}

type YearPage struct{ Page }

func (p YearPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.Year = atoiMust(p.Match[0])
	p.Articles = data.Articles.Year(p.Year)
	p.Title = fmt.Sprint(conf.Blog.Title, " - ", p.Year)
	p.Page.ServeHTTP(w, r)
}

type MonthPage struct{ Page }

func (p MonthPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.Year = atoiMust(p.Match[0])
	p.Month = time.Month(atoiMust(p.Match[1]))
	p.Articles = data.Articles.Year(p.Year).Month(p.Month)
	p.Title = fmt.Sprint(conf.Blog.Title, " - ", p.Year, p.Month)
	p.Page.ServeHTTP(w, r)
}
