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
	rss Rss
	sitemap SiteMap
)

func assetHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

/* temporary helper function */
func imgHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, conf.Blog.Url + r.URL.Path, http.StatusFound)
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
	rss = NewRss()
	sitemap = NewSitemap()

	http.HandleFunc("/assets/", assetHandler)
	http.HandleFunc("/images/", imgHandler)
	http.Handle("/tag/", TagPage{})
	/*
	http.HandleFunc("^/admin/(.+)$", adminSlug)
	http.HandleFunc("^/admin/?$", adminList)
	 */
	http.Handle("/rss.xml", rss)
	http.Handle("/sitemap.xml", sitemap)
	/*
	http.HandleFunc("^/(\\d+)/?$", year)
	http.HandleFunc("^/(\\d+)/(\\d+)/?$", month)
	 */
	http.Handle("/", IndexPage{})

	if err := http.ListenAndServe(listen, nil); err != nil {
		log.Fatal(err)
	}
}
