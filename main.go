package main

import (
	"flag"
	"log"
	"net/http"
	"sort"
	"text/template"

	"github.com/dim13/gold/articles"
	"github.com/dim13/gold/route"
	"github.com/dim13/gold/storage"
)

var (
	conf   storage.Config
	art    articles.Articles
	tmpl   *template.Template
	listen string
	config string
)

func assetHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

/* temporary helper function */
func tmpHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("redirect to dim13", r.URL.Path)
	http.Redirect(w, r, "http://www.dim13.org"+r.URL.Path, http.StatusFound)
}

func init() {
	flag.StringVar(&listen, "listen", ":8000", "listen at")
	flag.StringVar(&config, "config", "config/config.ini", "config file")
	flag.Parse()
}

func main() {
	var err error

	conf, err = storage.ReadConf(config)
	if err != nil {
		log.Fatal(err)
	}

	articles.SetStorage(conf.Blog.DataBase)
	err = art.Load()
	if err != nil {
		log.Println(err)
	}
	sort.Sort(sort.Reverse(art))

	log.Println("Prepare templates")
	tmpl = template.Must(template.ParseGlob("templates/*.tmpl"))

	re := new(route.ReHandler)
	re.HandleFunc("^/assets/", assetHandler)
	re.HandleFunc("^/(favicon\\.ico|images|videos)", tmpHandler)
	re.HandleFunc("^/robots.txt$", robotsHandler)
	re.HandleFunc("^/rss.xml$", rssHandler)
	re.HandleFunc("^/sitemap.xml$", sitemapHandler)
	re.Handle("^/admin/(.+)$", &adminSlug{})
	re.Handle("^/admin/?$", &adminIndex{})
	re.Handle("^/tags?/(.+)$", &tagPage{})
	re.Handle("^/\\d+/\\d+/(.+)$", &slugPage{})
	re.Handle("^/(\\d+)/(\\d+)/?$", &monthPage{})
	re.Handle("^/(\\d+)/?$", &yearPage{})
	re.Handle("^/(.+)$", &slugPage{})
	re.Handle("^/$", &indexPage{})

	log.Println("Listen on", listen)
	if err := http.ListenAndServe(listen, re); err != nil {
		log.Fatal(err)
	}
}
