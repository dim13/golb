// Test APP
package main

import (
	"github.com/dim13/gold"
	"log"
	"net/http"
	"sort"
	"text/template"
)

const listen = ":8000"

var (
	conf *gold.Config
	data *gold.Data
	tmpl *template.Template
)

type Page struct {
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
	tmpl.ExecuteTemplate(w, "admin.tmpl", p)
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

	tmpl.ExecuteTemplate(w, "admin.tmpl", p)
}

func index(w http.ResponseWriter, r *http.Request, s []string) {
	var a gold.Articles
	for _, v := range data.Articles {
		if v.Enabled {
			a = append(a, v)
		}
	}
	p := Page{
		Title:    conf.Blog.Title,
		Articles: a,
	}
	tmpl.ExecuteTemplate(w, "index.tmpl", p)
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
	tmpl.ExecuteTemplate(w, "index.tmpl", p)
}

func main() {
	var err error

	conf, err = gold.ReadConf("config/config.ini")
	if err != nil {
		log.Fatal(err)
	}

	data = gold.Open(conf.Settings.DataBase)
	if err := data.Read(); err != nil {
		log.Println(err)
	}
	sort.Sort(sort.Reverse(data.Articles))

	tmpl = template.Must(template.ParseGlob("templates/*.tmpl"))

	re := new(gold.ReHandler)
	re.AddRoute("^/admin/(.+)$", adminSlug)
	re.AddRoute("^/admin/?$", adminList)
	re.AddRoute("^/tags?/(.*)$", tags)
	re.AddRoute("^/(\\d+)/(\\d+)/(.*)$", index)
	re.AddRoute("^/(\\d+)/(.*)$", index)
	re.AddRoute("^/(.*)$", index)

	if err := http.ListenAndServe(listen, re); err != nil {
		log.Fatal(err)
	}
}
