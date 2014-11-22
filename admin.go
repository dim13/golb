package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dim13/gold/articles"
	"github.com/dim13/gold/storage"
)

type adminPage struct {
	Articles articles.Articles
	Title    string
	Config   storage.Config
	Error    string
}

func (p adminPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.Config = conf
	err := tmpl.ExecuteTemplate(w, "admin.tmpl", p)
	if err != nil {
		log.Println(err)
	}
}

type adminIndex struct{ adminPage }

func (p *adminIndex) Select(_ []string) bool {
	p.Articles = blog.Articles()
	p.Title = "Admin Interface"
	return true
}

func (p *adminIndex) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		p.adminPage.ServeHTTP(w, r)
		return
	}
	switch r.FormValue("submit") {
	case "add":
		a := articles.Article{
			Title: r.FormValue("title"),
		}
		r.URL.Path = "/admin/" + a.Slug()
	}
}

type adminSlug struct{ adminPage }

func (p *adminSlug) Select(match []string) bool {
	slug := match[0]

	if a, ok := blog[slug]; ok {
		p.Title = a.Title
		p.Articles = append(p.Articles, a)
	} else {
		p.Articles = append(p.Articles, articles.Article{
			Title: articles.MakeTitle(slug),
			Date:  time.Now(),
		})
	}
	return true
}

func (p *adminSlug) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		p.adminPage.ServeHTTP(w, r)
		return
	}
	a := articles.Article{
		Title:   r.FormValue("title"),
		Tags:    articles.ReadTags(r.FormValue("tags")),
		Body:    r.FormValue("body"),
		Enabled: r.FormValue("enabled") == "on",
	}
	/*
		if p.Articles[0].Slug() != a.Slug() {
			blog.Delete(p.Articles[0])
			r.URL.Path = "/admin/" + a.Slug()
		}
	*/
	switch r.FormValue("submit") {
	case "reload":
		log.Println("reloading")
		blog, _ = articles.Load()
	case "preview":
		blog.Add(a)
	case "save":
		blog.Add(a)
		blog.Store()
		r.URL.Path = "/admin/"
	case "delete":
		blog.Delete(a)
		blog.Store()
		r.URL.Path = "/admin/"
	case "cancel":
		r.URL.Path = "/admin/"
	}
}
