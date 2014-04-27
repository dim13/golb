// Test APP
package main

import (
	"github.com/dim13/gold"
	"log"
	"net/http"
	"sort"
	"strconv"
	"text/template"
	"time"
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
	PrevPage int
	NextPage int
	Expand   bool
	TagCloud gold.TagCloud
}

func adminList(w http.ResponseWriter, r *http.Request, s []string) {
	p := Page{
		Config:   conf,
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
		p = Page{
			Config: conf,
			Error:  err,
		}
	} else {
		p = Page{
			Config:  conf,
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
	var p Page

	a, err := data.Articles.Find(s[0])
	if err == nil {
		p = Page{
			Title:    a.Title,
			Articles: gold.Articles{a},
			Expand:   true,
		}
		genpage(w, p)
	} else {
		page(w, r, []string{"1"})
	}
}

func page(w http.ResponseWriter, r *http.Request, s []string) {
	pg, err := strconv.Atoi(s[0])
	if err != nil {
		log.Fatal(err)
	}
	app := conf.Blog.ArticlesPerPage

	a, next, prev := data.Articles.Enabled().Page(pg, app)

	p := Page{
		Title:    conf.Blog.Title,
		Articles: a,
		NextPage: next,
		PrevPage: prev,
	}

	genpage(w, p)
}

func tags(w http.ResponseWriter, r *http.Request, s []string) {
	var a gold.Articles
	for _, v := range data.Articles {
		if v.Tags.Has(s[0]) {
			a = append(a, v)
		}
	}
	p := Page{
		Title:    conf.Blog.Title + " - " + s[0],
		Articles: a,
	}
	genpage(w, p)
}

func assets(w http.ResponseWriter, r *http.Request, s []string) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func rss(w http.ResponseWriter, r *http.Request, s []string) {
	a := data.Articles.Enabled()
	app := conf.Blog.ArticlesPerPage

	p := Page{
		Config:   conf,
		Articles: a[:app],
	}

	err := tmpl.ExecuteTemplate(w, "rss.tmpl", p)
	if err != nil {
		log.Fatal(err)
	}
}

type Sitemap struct {
	Loc        string
	Date       time.Time
	ChangeFreq string
	Priority   float64
}

func (s Sitemap) LastMod() string {
	return s.Date.Local().Format("2006-02-01")
}

func sitemap(w http.ResponseWriter, r *http.Request, s []string) {
	var sm []Sitemap
	sm = append(sm, Sitemap{
		Loc:      conf.Blog.Url,
		Priority: 1.0,
	})
	for _, a := range data.Articles.Enabled() {
		sm = append(sm, Sitemap{
			Loc:        conf.Blog.Url + "/" + a.Slug,
			Priority:   0.8,
			Date:       a.Date,
			ChangeFreq: "monthly",
		})
	}
	for _, t := range data.Articles.TagCloud() {
		sm = append(sm, Sitemap{
			Loc:      conf.Blog.Url + "/tag/" + t.Tag,
			Priority: 0.6 - float64(t.Wight)/10,
		})
	}
	err := tmpl.ExecuteTemplate(w, "sitemap.tmpl", sm)
	if err != nil {
		log.Fatal(err)
	}
}

func year(w http.ResponseWriter, r *http.Request, s []string) {
	y, err := strconv.Atoi(s[0])
	if err != nil {
		log.Fatal(err)
	}
	a := data.Articles.Year(y)
	p := Page{
		Title:    conf.Blog.Title + " - " + s[0],
		Articles: a,
	}
	genpage(w, p)
}

func genpage(w http.ResponseWriter, p Page) {
	tic := conf.Blog.TagsInCloud
	p.TagCloud = data.Articles.TagCloud()[:tic]
	p.Config = conf

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
	re.AddRoute("^/assets/", assets)
	re.AddRoute("^/images/", assets)
	re.AddRoute("^/admin/(.+)$", adminSlug)
	re.AddRoute("^/admin/?$", adminList)
	re.AddRoute("^/tags?/(.+)$", tags)
	re.AddRoute("^/page/(\\d+)$", page)
	re.AddRoute("^/rss\\.xml$", rss)
	re.AddRoute("^/sitemap\\.xml$", sitemap)
	re.AddRoute("^/(\\d+)/$", year)
	re.AddRoute("^/(\\d+)/(\\d+)/(.*)$", index)
	re.AddRoute("^/(.*)$", index)

	if err := http.ListenAndServe(listen, re); err != nil {
		log.Fatal(err)
	}
}
