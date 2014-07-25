package main

import (
	"log"
	"net/http"

	"github.com/dim13/gold/articles"
	"github.com/dim13/gold/storage"
)

type AdminPage struct {
	Articles articles.Articles
	Article  *articles.Article
	Title    string
	Config   storage.Config
	Error    string
}

func (p AdminPage) Get(w http.ResponseWriter, r *http.Request) {
	p.Config = conf
	err := tmpl.ExecuteTemplate(w, "admin.tmpl", p)
	if err != nil {
		log.Println(err)
	}
}

func (p *AdminPage) Post(w http.ResponseWriter, r *http.Request) {
	log.Println("Catch POST redirect admin", r.URL.Path)
}

type AdminIndex struct{ AdminPage }

func (p *AdminIndex) Select(_ []string) {
	p.Articles = art
	p.Title = "Admin Interface"
}

func (p *AdminIndex) Post(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	r.URL.Path = "/admin/" + articles.MakeSlug(title)
}

type AdminSlug struct{ AdminPage }

func (p *AdminSlug) Select(match []string) {
	if a, ok := art.Find(match[0]); ok {
		p.Title = a.Title
		p.Article = a
	} else {
		p.Article = &articles.Article{
			Slug: match[0],
			Title: articles.MakeTitle(match[0]),
		}
	}
}

func (p *AdminSlug) Post(w http.ResponseWriter, r *http.Request) {
	a := articles.Article{
		Title:   r.FormValue("title"),
		Slug:    r.FormValue("slug"),
		Tags:    articles.ReadTags(r.FormValue("tags")),
		Body:    r.FormValue("body"),
		Enabled: r.FormValue("enabled") == "on",
	}
	art.Update(a)
	r.URL.Path = "/admin/"
}
