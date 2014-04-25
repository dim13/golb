// Test APP
package main

import (
	"github.com/dim13/gold"
	"log"
	"net/http"
	"sort"
	"strconv"
	"text/template"
)

const listen = ":8000"

var (
	conf *gold.Config
	data *gold.Data
	tmpl *template.Template
)

type Page struct {
	Config   *gold.Config
	Title    string
	Articles gold.Articles
	Article  *gold.Article
	Error    error
	PrevPage int
	NextPage int
	Expand   bool
	TagCloud gold.TagCloud
}

func adminList(w http.ResponseWriter, r *http.Request, s []string) {
	p := Page{
		Config:   conf,
		Title:    "Admin interface",
		Articles: data.Articles,
	}
	err := tmpl.ExecuteTemplate(w, "admin.tmpl", p)
	if err != nil {
		log.Fatal(err)
	}
}

func adminSlug(w http.ResponseWriter, r *http.Request, s []string) {
	var p Page

	a, err := data.Articles.Find(s[0])
	if err != nil {
		p = Page{
			Config: conf,
			Error:  err,
		}
	} else {
		p = Page{
			Config:  conf,
			Title:   a.Title,
			Article: a,
		}
	}

	err = tmpl.ExecuteTemplate(w, "admin.tmpl", p)
	if err != nil {
		log.Fatal(err)
	}
}

func index(w http.ResponseWriter, r *http.Request, s []string) {
	var p Page

	a, err := data.Articles.Find(s[0])
	if err == nil {
		p = Page{
			Config:   conf,
			Title:    a.Title,
			Articles: gold.Articles{a},
			Expand:   true,
		}
		err = tmpl.ExecuteTemplate(w, "index.tmpl", p)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		page(w, r, []string{"1"})
	}
}

func page(w http.ResponseWriter, r *http.Request, s []string) {
	var (
		a    gold.Articles
		prev int
		next int
	)

	tc := data.Articles.TagCloud(conf.Blog.TagsInCloud)
	n := conf.Blog.ArticlesPerPage

	pg, err := strconv.Atoi(s[0])
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range data.Articles {
		if v.Enabled {
			a = append(a, v)
		}
	}

	/* sanitize lower bound */
	if pg < 1 {
		pg = 1
	}

	/* sanitize upper bound */
	m := len(a)/n + 1
	if pg > m {
		pg = m
	}

	first := (pg - 1) * n
	last := first + n - 1
	if last > len(a) {
		last = len(a)
	}

	if pg > 1 {
		prev = pg - 1
	}

	if pg < m {
		next = pg + 1
	}

	p := Page{
		Config:   conf,
		Title:    conf.Blog.Title,
		Articles: a[first:last],
		NextPage: next,
		PrevPage: prev,
		TagCloud: tc,
	}

	err = tmpl.ExecuteTemplate(w, "index.tmpl", p)
	if err != nil {
		log.Fatal(err)
	}
}

func tags(w http.ResponseWriter, r *http.Request, s []string) {
	var a gold.Articles
	for _, v := range data.Articles {
		if v.Tags.Has(s[0]) {
			a = append(a, v)
		}
	}
	p := Page{
		Config:   conf,
		Title:    conf.Blog.Title + " - " + s[0],
		Articles: a,
	}
	err := tmpl.ExecuteTemplate(w, "index.tmpl", p)
	if err != nil {
		log.Fatal(err)
	}
}

func assets(w http.ResponseWriter, r *http.Request, s []string) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func rss(w http.ResponseWriter, r *http.Request, s []string) {
	var a gold.Articles
	for i, v := range data.Articles {
		if v.Enabled {
			a = append(a, v)
		}
		if i > conf.Blog.ArticlesPerPage {
			break
		}
	}
	p := Page{
		Config: conf,
		Articles: a,
	}
	err := tmpl.ExecuteTemplate(w, "rss.tmpl", p)
	if err != nil {
		log.Fatal(err)
	}
}

func year(w http.ResponseWriter, r *http.Request, s []string) {
	y, err := strconv.Atoi(s[0])
	if err != nil {
		log.Fatal(err)
	}
	a := data.Articles.Year(y)
	p := Page{
		Config:   conf,
		Title:    conf.Blog.Title + " - " + s[0],
		Articles: a,
	}
	err = tmpl.ExecuteTemplate(w, "index.tmpl", p)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var err error

	conf, err = gold.ReadConf("config/config.ini")
	if err != nil {
		log.Fatal(err)
	}

	data = gold.Open(conf.Blog.DataBase)
	if err := data.Read(); err != nil {
		log.Println(err)
	}
	sort.Sort(sort.Reverse(data.Articles))

	tmpl = template.Must(template.ParseGlob("templates/*.tmpl"))

	re := new(gold.ReHandler)
	re.AddRoute("^/assets/", assets)
	re.AddRoute("^/images/", assets)
	re.AddRoute("^/admin/(.+)$", adminSlug)
	re.AddRoute("^/admin/?$", adminList)
	re.AddRoute("^/tags?/(.+)$", tags)
	re.AddRoute("^/page/(\\d+)$", page)
	re.AddRoute("^/rss\\.xml$", rss)
	re.AddRoute("^/(\\d+)/$", year)
	re.AddRoute("^/(\\d+)/(\\d+)/(.*)$", index)
	re.AddRoute("^/(.*)$", index)

	if err := http.ListenAndServe(listen, re); err != nil {
		log.Fatal(err)
	}
}
