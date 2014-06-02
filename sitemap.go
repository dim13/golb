package main

import (
	"log"
	"net/http"
	"time"
)

type SiteMap []Sitemap
type Sitemap struct {
	Loc        string
	Date       time.Time
	ChangeFreq string
	Priority   float64
}

func (sm Sitemap) LastMod() string {
	if sm.Date.IsZero() {
		return ""
	}
	return sm.Date.Local().Format("2006-02-01")
}

func sitemapHandler(w http.ResponseWriter, r *http.Request) {
	var sm SiteMap
	articles := data.Articles.Enabled()
	sm = append(sm, Sitemap{
		Loc:        "http://" + r.Host,
		Priority:   1.0,
		Date:       articles[0].Date,
		ChangeFreq: "daily",
	})
	for _, a := range articles {
		sm = append(sm, Sitemap{
			Loc:        "http://" + r.Host + "/" + a.Slug,
			Priority:   0.8,
			Date:       a.Date,
			ChangeFreq: "monthly",
		})
	}
	for _, t := range data.Articles.TagCloud() {
		tagged := articles.Tag(t.Tag)
		sm = append(sm, Sitemap{
			Loc:        "http://" + r.Host + "/tag/" + t.Tag,
			Priority:   0.6 - float64(t.Wight)/10,
			Date:       tagged[0].Date,
			ChangeFreq: "weekly",
		})
	}
	err := tmpl.ExecuteTemplate(w, "sitemap.tmpl", sm)
	if err != nil {
		log.Fatal(err)
	}
}
