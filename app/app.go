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
	re.AddRoute("^/assets/", assets)
	re.AddRoute("^/images/", assets)
	re.AddRoute("^/admin/(.+)$", adminSlug)
	re.AddRoute("^/admin/?$", adminList)
	re.AddRoute("^/tags?/(.+)$", tags)
	re.AddRoute("^/page/(\\d+)$", page)
	re.AddRoute("^/rss\\.xml$", rss)
	re.AddRoute("^/sitemap\\.xml$", sitemap)
	re.AddRoute("^/(\\d+)/?$", year)
	re.AddRoute("^/(\\d+)/(\\d+)/?$", month)
	re.AddRoute("^/\\d+/\\d+/(.*)$", redir)
	re.AddRoute("^/(.*)$", index)

	if err := http.ListenAndServe(listen, re); err != nil {
		log.Fatal(err)
	}
}
