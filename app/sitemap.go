package main

import (
	"log"
	"net/http"
	"time"
)

type Sitemap struct {
	Loc        string
	Date       time.Time
	ChangeFreq string
	Priority   float64
}

func (s Sitemap) LastMod() string {
	return s.Date.Local().Format("2006-02-01")
}

func sitemapHandler(w http.ResponseWriter, r *http.Request, s []string) {
	var sm []Sitemap
	sm = append(sm, Sitemap{
		Loc:      conf.Blog.Url,
		Priority: 1.0,
	})
	for _, a := range data.Articles.Enabled() {
		sm = append(sm, Sitemap{
			Loc:        conf.Blog.Url + "/" + a.Slug,
			Priority:   0.8,
			Date:       a.Date,
			ChangeFreq: "monthly",
		})
	}
	for _, t := range data.Articles.TagCloud() {
		sm = append(sm, Sitemap{
			Loc:      conf.Blog.Url + "/tag/" + t.Tag,
			Priority: 0.6 - float64(t.Wight)/10,
		})
	}
	err := tmpl.ExecuteTemplate(w, "sitemap.tmpl", sm)
	if err != nil {
		log.Fatal(err)
	}
}
