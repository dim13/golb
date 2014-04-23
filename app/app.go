// Test APP
package main

import (
	"fmt"
	"github.com/dim13/gold"
	"text/template"
	"log"
	"net/http"
	"sort"
)

const listen = ":8000"

var (
	conf *gold.Config
	data *gold.Data
	tmpl *template.Template
)

type Admin struct {
	Title    string
	Articles gold.Articles
	Article  *gold.Article
	Error    error
}

func adminList(w http.ResponseWriter, r *http.Request, s []string) {
	p := Admin{
		Title:    "Admin interface",
		Articles: data.Articles,
	}
	tmpl.ExecuteTemplate(w, "admin.tmpl", p)
}

func adminSlug(w http.ResponseWriter, r *http.Request, s []string) {
	var p Admin

	a, err := data.Articles.Find(s[0])
	if err != nil {
		p = Admin{Error: err}
	} else {
		p = Admin{
			Title:   a.Title,
			Article: a,
		}
	}

	tmpl.ExecuteTemplate(w, "admin.tmpl", p)
}

func root(w http.ResponseWriter, r *http.Request, s []string) {
	fmt.Fprint(w, s)
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
	re.AddRoute("^/(\\d+)/(\\d+)/(.*)$", root)
	re.AddRoute("^/(\\d+)/(.*)$", root)
	re.AddRoute("^/(.*)$", root)

	if err := http.ListenAndServe(listen, re); err != nil {
		log.Fatal(err)
	}
}
