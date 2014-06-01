package main

import (
	"github.com/dim13/gold/storage"
	"log"
	"net/http"
	"sort"
	"text/template"
)

const listen = ":8000"

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
	http.Redirect(w, r, conf.Blog.Url + r.URL.Path, http.StatusFound)
}

func main() {
	var err error

	conf, err = storage.ReadConf("config/config.ini")
	if err != nil {
		log.Fatal(err)
	}

	data = storage.Open(conf.Blog.DataBase)
	if err := data.Read(); err != nil {
		log.Println(err)
	}
	sort.Sort(sort.Reverse(data.Articles))

	tmpl = template.Must(template.ParseGlob("templates/*.tmpl"))

	re := new(ReHandler)

	re.HandleFunc("^/assets/", assetHandler)
	re.HandleFunc("^/images/", tmpHandler)
	re.HandleFunc("^/videos/", tmpHandler)
	re.Handle("^/rss.xml$", &Rss{})
	re.Handle("^/sitemap.xml$", &SiteMap{})
	re.Handle("^/admin/(.+)$", &AdminSlug{})
	re.Handle("^/admin/?$", &AdminIndex{})
	re.Handle("^/tags?/(.+)$", &TagPage{})
	re.Handle("^/\\d+/\\d+/(.+)$", &SlugPage{})
	re.Handle("^/(\\d+)/(\\d+)/?$", &MonthPage{})
	re.Handle("^/(\\d+)/?$", &YearPage{})
	re.Handle("^/(.+)$", &SlugPage{})
	re.Handle("^/$", &IndexPage{})

	if err := http.ListenAndServe(listen, re); err != nil {
		log.Fatal(err)
	}
}
