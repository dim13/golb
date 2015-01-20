package main

import (
	"log"
	"net/http"
	//"time"

	"github.com/dim13/gold/blog"
	"github.com/dim13/gold/storage"
)

type adminPage struct {
	Articles blog.Articles
	Article  *blog.Article
	Title    string
	Config   storage.Config
	Error    string
}

func (p adminPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.Config = Conf
	err := tmpl.ExecuteTemplate(w, "admin.tmpl", p)
	if err != nil {
		log.Println(err)
	}
}

func adminIndexHandler(w http.ResponseWriter, r *http.Request) {
	pg := adminPage{
		Articles: Blog.Articles(),
		Title:    "Published Articles",
	}
	pg.ServeHTTP(w, r)
}

func adminDraftHandler(w http.ResponseWriter, r *http.Request) {
	pg := adminPage{
		Articles: Blog.Drafts(),
		Title:    "Unpublished Articles",
	}
	pg.ServeHTTP(w, r)
}

/*
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
*/

func adminSlugHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get(":slug")
	a, ok := Blog.Public[slug]
	if !ok {
		a, _ = Blog.Draft[slug]
	}
	pg := adminPage{
		Title:   a.Title,
		Article: &a,
	}
	pg.ServeHTTP(w, r)
}

/*
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
*/
