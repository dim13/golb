package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dim13/gold/blog"
	"github.com/dim13/gold/storage"
)

type page struct {
	Config    storage.Config
	URL       string
	Title     string
	Articles  blog.Articles
	Error     error
	PrevPage  int
	NextPage  int
	TagCloud  blog.TagCloud
	Year      int
	Month     time.Month
	Archive   []year
	FirstYear int
	LastYear  int
}

func atoiMust(s string) int {
	if s == "" {
		return 0
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Println(err)
		return 0
	}
	return i
}

func (p *page) Pager(pg string) {
	perpage := Conf.Blog.ArticlesPerPage
	count := len(p.Articles)
	last := count/perpage + 1
	curr := atoiMust(pg)

	if curr <= 1 {
		curr = 1
	} else {
		p.PrevPage = curr - 1
	}

	if curr >= last {
		curr = last
	} else {
		p.NextPage = curr + 1
	}

	from := (curr - 1) * perpage
	to := from + perpage - 1

	if to > count {
		to = count
	}

	p.Articles = p.Articles[from:to]
}

func (p page) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.URL = "http://" + r.Host
	p.Pager(r.URL.Query().Get("page"))
	p.TagCloud = Blog.TagCloud()
	p.Config = Conf
	if p.Year == 0 {
		p.Year = p.Articles.Head().Year()
	}
	if p.Month == 0 {
		p.Month = p.Articles.Head().Month()
	}
	p.MakeArchive()

	a := Blog.Articles()
	p.FirstYear = a.Tail().Year()
	p.LastYear = a.Head().Year()

	err := tmpl.ExecuteTemplate(w, "index.tmpl", p)
	if err != nil {
		log.Println(err)
	}
}

func tagHandler(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get(":tag")
	tagged, _ := Blog.TagMap()[tag]
	pg := page{
		Articles: tagged,
		Title:    tag,
	}
	pg.ServeHTTP(w, r)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	pg := page{
		Articles: Blog.Articles(),
	}
	pg.ServeHTTP(w, r)
}

func slugHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get(":slug")
	a, _ := Blog.Public[slug]
	pg := page{
		Title:    a.Title,
		Articles: blog.Articles{a},
		Year:     a.Year(),
		Month:    a.Month(),
	}
	pg.ServeHTTP(w, r)
}

func yearHandler(w http.ResponseWriter, r *http.Request) {
	year := atoiMust(r.URL.Query().Get(":year"))
	pg := page{
		Year:     year,
		Articles: Blog.Articles().Year(year),
		Title:    fmt.Sprint(year),
	}
	pg.ServeHTTP(w, r)
}

func monthHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	year := atoiMust(q.Get(":year"))
	month := time.Month(atoiMust(q.Get(":month")))
	pg := page{
		Year:     year,
		Month:    month,
		Articles: Blog.Articles().Year(year).Month(month),
		Title:    fmt.Sprint(month, year),
	}
	pg.ServeHTTP(w, r)
}
