package main

import (
	"fmt"
	"github.com/dim13/gold/storage"
	"log"
	"net/http"
	"sort"
	"text/template"
)

const (
	listen = ":8000"
	config = "config/config.ini"
)

var (
	conf *storage.Config
	data *storage.Data
	tmpl *template.Template
)

func assetHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

/* temporary helper function */
func tmpHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://"+r.Host+r.URL.Path, http.StatusFound)
}

func robotsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "User-agent: *")
	fmt.Fprintln(w, "Sitemap:", "http://"+r.Host+"/sitemap.xml")
}

func main() {
	var err error

	log.Println("Read", config)
	conf, err = storage.ReadConf(config)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Open", conf.Blog.DataBase)
	data = storage.Open(conf.Blog.DataBase)
	if err := data.Read(); err != nil {
		log.Println(err)
	}
	sort.Sort(sort.Reverse(data.Articles))

	log.Println("Prepare templates")
	tmpl = template.Must(template.ParseGlob("templates/*.tmpl"))

	re := new(ReHandler)
	re.HandleFunc("^/assets/", assetHandler)
	re.HandleFunc("^/(images|videos)/", tmpHandler)
	re.HandleFunc("^/robots.txt$", robotsHandler)
	re.HandleFunc("^/rss.xml$", rssHandler)
	re.HandleFunc("^/sitemap.xml$", sitemapHandler)
	re.Handle("^/admin/(.+)$", &AdminSlug{})
	re.Handle("^/admin/?$", &AdminIndex{})
	re.Handle("^/tags?/(.+)$", &TagPage{})
	re.Handle("^/\\d+/\\d+/(.+)$", &SlugPage{})
	re.Handle("^/(\\d+)/(\\d+)/?$", &MonthPage{})
	re.Handle("^/(\\d+)/?$", &YearPage{})
	re.Handle("^/(.+)$", &SlugPage{})
	re.Handle("^/$", &IndexPage{})

	log.Println("Listen on", listen)
	if err := http.ListenAndServe(listen, re); err != nil {
		log.Fatal(err)
	}
}
