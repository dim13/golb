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

func redirHandler(w http.ResponseWriter, r *http.Request, s []string) {
	http.Redirect(w, r, "/" + s[0], http.StatusFound)
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
	re.AddRoute("^/assets/", assetsHandler)
	re.AddRoute("^/images/", assetsHandler)
	re.AddRoute("^/admin/(.+)$", adminSlugHandler)
	re.AddRoute("^/admin/?$", adminListHandler)
	re.AddRoute("^/tags?/(.+)$", tagsHandler)
	re.AddRoute("^/page/(\\d+)$", pageHandler)
	re.AddRoute("^/rss\\.xml$", rssHandler)
	re.AddRoute("^/sitemap\\.xml$", sitemapHandler)
	re.AddRoute("^/(\\d+)/?$", yearHandler)
	re.AddRoute("^/(\\d+)/(\\d+)/?$", monthHandler)
	re.AddRoute("^/\\d+/\\d+/(.*)$", redirHandler)
	re.AddRoute("^/(.*)$", indexHandler)

	if err := http.ListenAndServe(listen, re); err != nil {
		log.Fatal(err)
	}
}
