package main

import (
	"github.com/dim13/gold"
	"log"
	"net/http"
)

type AdminPage struct {
	Articles gold.Articles
	Article  *gold.Article
	Title    string
	Config   *gold.Config
	Error    string
}

func (p AdminPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.Config = conf
	err := tmpl.ExecuteTemplate(w, "admin.tmpl", p)
	if err != nil {
		log.Fatal(err)
	}
}

type AdminIndex struct { AdminPage }

func (p *AdminIndex) Select(match []string) {
	p.Articles = data.Articles
	p.Title = "Admin Interface"
}

func (p *AdminIndex) Store(r *http.Request) {
	log.Println(p, r)
}

type AdminSlug struct { AdminPage }

func (p *AdminSlug) Select(match []string) {
	a, err := data.Articles.Find(match[0])
	if err == nil {
		p.Title = a.Title
		p.Article = a
	}
}

func (p *AdminSlug) Store(r *http.Request) {
	a := gold.Article{
		Title: r.FormValue("title"),
		Slug: r.FormValue("slug"),
		Tags: gold.ReadTags(r.FormValue("tags")),
		Body: r.FormValue("body"),
		Enabled: r.FormValue("enabled") != "",
	}
	p.Article = &a
	if r.FormValue("save") != "" {
		data.Articles.Update(&a)
	}
	//log.Println(p, r)
}
