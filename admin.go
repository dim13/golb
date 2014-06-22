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
	Config   *storage.Config
	Error    string
}

func (p AdminPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.Config = conf
	err := tmpl.ExecuteTemplate(w, "admin.tmpl", p)
	if err != nil {
		log.Println(err)
	}
}

type AdminIndex struct{ AdminPage }

func (p *AdminIndex) Select(match []string) {
	p.Articles = data.Articles
	p.Title = "Admin Interface"
}

func (p *AdminIndex) Store(r *http.Request) {
	//log.Println(p, r)
}

type AdminSlug struct{ AdminPage }

func (p *AdminSlug) Select(match []string) {
	a, err := data.Articles.Find(match[0])
	if err == nil {
		p.Title = a.Title
		p.Article = a
	}
}

func (p *AdminSlug) Store(r *http.Request) {
	a := articles.Article{
		Title:   r.FormValue("title"),
		Slug:    r.FormValue("slug"),
		Tags:    articles.ReadTags(r.FormValue("tags")),
		Body:    r.FormValue("body"),
		Enabled: r.FormValue("enabled") != "",
	}
	p.Article = &a
	if r.FormValue("save") != "" {
		data.Articles.Update(&a)
	}
	//log.Println(p, r)
}
