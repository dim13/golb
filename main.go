package main

import (
	"flag"
	"log"
	"net/http"
	"text/template"

	"github.com/dim13/gold/blog"
	"github.com/dim13/gold/storage"

	"github.com/bmizerany/pat"
)

var (
	Conf storage.Config
	Blog blog.Blog
	tmpl *template.Template

	listen string
	config string
)

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	err := tmpl.ExecuteTemplate(w, "notfound.tmpl", nil)
	if err != nil {
		log.Println(err)
	}
}

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

	Conf, err = storage.ReadConf(config)
	if err != nil {
		log.Fatal(err)
	}

	blog.SetStorage(Conf.Blog.DataBase)
	Blog, err = blog.Load()
	if err != nil {
		log.Println(err)
	}

	log.Println("Prepare templates")
	tmpl = template.Must(template.ParseGlob("templates/*.tmpl"))

	mux := pat.New()

	mux.Get("/assets/", http.HandlerFunc(assetHandler))
	mux.Get("/images/", http.HandlerFunc(tmpHandler))
	mux.Get("/videos/", http.HandlerFunc(tmpHandler))

	mux.Get("/robots.txt", robotsPage{})
	mux.Get("/sitemap.xml", sitemapPage{})
	mux.Get("/rss.xml", rssPage{})

	mux.Get("/admin/drafts/", http.HandlerFunc(adminDraftHandler))
	mux.Get("/admin/:slug", http.HandlerFunc(adminSlugHandler))
	mux.Get("/admin/", http.HandlerFunc(adminIndexHandler))

	mux.Get("/tag/:tag", http.HandlerFunc(tagHandler))
	mux.Get("/:year/:month/:slug", http.HandlerFunc(slugHandler))
	mux.Get("/:year/:month/", http.HandlerFunc(monthHandler))
	mux.Get("/:year/", http.HandlerFunc(yearHandler))
	mux.Get("/", http.HandlerFunc(indexHandler))

	log.Println("Listen on", listen)
	log.Fatal(http.ListenAndServe(listen, mux))
}
