// Test APP
package main

import (
	"fmt"
	"github.com/dim13/gold"
	"log"
	"net/http"
	"html/template"
)

const listen = ":8000"

var (
	conf *gold.Config
	data *gold.Data
	tmpl *template.Template
)

type Admin struct {
	Title string
	Articles gold.Articles
}

func admin(w http.ResponseWriter, r *http.Request, s []string) {
	p := &Admin {
		Title: "Admin interface",
		Articles: data.Articles,
	}
	tmpl.Execute(w, p)
}

func adminAdd(w http.ResponseWriter, r *http.Request, s []string) {
	fmt.Fprint(w, s)
}

func adminSlug(w http.ResponseWriter, r *http.Request, s []string) {
	fmt.Fprint(w, s)
}

func root(w http.ResponseWriter, r *http.Request, s []string) {
	fmt.Fprint(w, s)
}

func main() {
	var err error

	conf, err = gold.ReadConf("config.ini")
	if err != nil {
		log.Fatal(err)
	}

	data = gold.Open(conf.Settings.DataBase)
	if err := data.Read(); err != nil {
		log.Println(err)
	}

	tmpl = template.Must(template.ParseFiles("admin.tmpl"))

	a := &gold.Article{
		Title: "Test title",
		Body: "empty body",
		Tags: []string{"no", "tags", "at all"},
		Author: "me@example.com",
	}
	err = data.Articles.Add(a)
	if err != nil {
		log.Println(err)
	}

	c := &gold.Comment{
		Name: "anonymous coward",
		Email: "none@example.com",
		URL: "http://example.com",
		Comment: "empty",
	}
	a.Comments.Add(c)
	c.Publish()

	z, err := data.Articles.Find("test-title")
	if err == nil {
		z.Publish()
	}

	if err := data.Write(); err != nil {
		log.Fatal(err)
	}


	re := new(gold.ReHandler)
	re.AddRoute("^/admin/?$", admin)
	re.AddRoute("^/admin/add$", adminAdd)
	re.AddRoute("^/admin/(.*)$", adminSlug)
	re.AddRoute("^/(\\d+)/(.*)$", root)
	if err := http.ListenAndServe(listen, re); err != nil {
		log.Fatal(err)
	}
}
