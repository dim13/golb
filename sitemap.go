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

type sitemapPage struct {
	Sitemap []sitemap
}

func (p sitemapPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.Sitemap = append(p.Sitemap, sitemap{
		Loc:      "http://" + r.Host,
		Priority: 1.0,
	})

	for _, a := range Blog.Articles() {
		p.Sitemap = append(p.Sitemap, sitemap{
			Loc:      "http://" + r.Host + "/" + a.Slug,
			Priority: 0.8,
		})
	}

	for t, a := range Blog.TagMap() {
		p.Sitemap = append(p.Sitemap, sitemap{
			Loc:      "http://" + r.Host + "/tag/" + t,
			Priority: 0.6 - float64(5/len(a))/10,
		})
	}

	err := tmpl.ExecuteTemplate(w, "sitemap.tmpl", p)
	if err != nil {
		log.Println(err)
	}
}
