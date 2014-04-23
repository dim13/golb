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
}

func adminList(w http.ResponseWriter, r *http.Request, s []string) {
	p := Page{
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
		p = Page{Error: err}
	} else {
		p = Page{
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
	page(w, r, []string{"0"})
}

func page(w http.ResponseWriter, r *http.Request, s []string) {
	var a gold.Articles
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

	first := (pg * n) % len(a)
	last := first + n
	if last > len(a) {
		last = len(a)
	}

	p := Page{
		Config:   conf,
		Title:    conf.Blog.Title,
		Articles: a[first : last],
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
		Title:    s[0],
		Articles: a,
	}
	err := tmpl.ExecuteTemplate(w, "index.tmpl", p)
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
	re.AddRoute("^/admin/(.+)$", adminSlug)
	re.AddRoute("^/admin/?$", adminList)
	re.AddRoute("^/tags?/(.+)$", tags)
	re.AddRoute("^/page/(\\d+)$", page)
	re.AddRoute("^/(\\d+)/(\\d+)/(.*)$", index)
	re.AddRoute("^/(\\d+)/(.*)$", index)
	re.AddRoute("^/(.*)$", index)

	if err := http.ListenAndServe(listen, re); err != nil {
		log.Fatal(err)
	}
}
