// Test APP
package main

import (
	"fmt"
	"github.com/dim13/gold"
	"html/template"
	"log"
	"net/http"
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

func admin(w http.ResponseWriter, r *http.Request, s []string) {
	log.Println(s)
	var p Admin
	if len(s) == 2 {
		a, err := data.Articles.Find(s[1])
		if err != nil {
			p = Admin{Error: err}
		} else {
			p = Admin{
				Title:   a.Title,
				Article: a,
			}
		}
	} else {
		p = Admin{
			Title:    "Admin interface",
			Articles: data.Articles,
		}
	}
	log.Println(p.Article)
	err := tmpl.ExecuteTemplate(w, "admin.tmpl", p)
	if err != nil {
		log.Println(err)
	}
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

	tmpl = template.Must(template.ParseGlob("templates/*.tmpl"))

	re := new(gold.ReHandler)
	re.AddRoute("^/admin/(.*)$", admin)
	re.AddRoute("^/(\\d+)/(.*)$", root)
	re.AddRoute("^/(\\d+)/(\\d+)(.*)$", root)
	re.AddRoute("^/(.*)$", root)
	if err := http.ListenAndServe(listen, re); err != nil {
		log.Fatal(err)
	}
}
