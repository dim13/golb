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
)

func redir(w http.ResponseWriter, r *http.Request, s []string) {
	http.Redirect(w, r, "/"+s[0], http.StatusFound)
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

	http.HandleFunc("/assets/", assetHandler)
	http.HandleFunc("/images/", assetHandler)
	/*
	http.HandleFunc("^/admin/(.+)$", adminSlug)
	http.HandleFunc("^/admin/?$", adminList)
	 */
	http.HandleFunc("/rss.xml", rssHandler)
	http.HandleFunc("/sitemap.xml", sitemapHandler)
	/*
	http.HandleFunc("^/(\\d+)/?$", year)
	http.HandleFunc("^/(\\d+)/(\\d+)/?$", month)
	http.HandleFunc("^/\\d+/\\d+/(.*)$", redir)
	 */
	http.HandleFunc("/", indexHandler)

	if err := http.ListenAndServe(listen, nil); err != nil {
		log.Fatal(err)
	}
}
