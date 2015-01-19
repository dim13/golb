package main

import (
	"log"
	"net/http"
)

type sitemap struct {
	Loc        string
	LastMod    string
	ChangeFreq string
	Priority   float64
}

func sitemapHandler(w http.ResponseWriter, r *http.Request) {
	var sm []sitemap

	sm = append(sm, sitemap{
		Loc:      "http://" + r.Host,
		Priority: 1.0,
	})

	for _, a := range Blog.Articles() {
		sm = append(sm, sitemap{
			Loc:      "http://" + r.Host + "/" + a.Slug,
			Priority: 0.8,
		})
	}

	for t, a := range Blog.TagMap() {
		sm = append(sm, sitemap{
			Loc:      "http://" + r.Host + "/tag/" + t,
			Priority: 0.6 - float64(5/len(a))/10,
		})
	}

	err := tmpl.ExecuteTemplate(w, "sitemap.tmpl", sm)
	if err != nil {
		log.Println(err)
	}
}
