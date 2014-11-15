package main

import (
	"log"
	"net/http"
	"time"
)

type sitemap struct {
	Loc        string
	LastMod    string
	ChangeFreq string
	Priority   float64
}

func lastMod(date time.Time) string {
	if date.IsZero() {
		return ""
	}
	return date.Local().Format("2006-02-01")
}

func sitemapHandler(w http.ResponseWriter, r *http.Request) {
	var sm []sitemap
	articles := art.Enabled()
	sm = append(sm, sitemap{
		Loc:        "http://" + r.Host,
		Priority:   1.0,
		LastMod:    lastMod(articles[0].Date),
		ChangeFreq: "daily",
	})
	for _, a := range articles {
		sm = append(sm, sitemap{
			Loc:        "http://" + r.Host + "/" + a.Slug,
			Priority:   0.8,
			LastMod:    lastMod(a.Date),
			ChangeFreq: "monthly",
		})
	}
	for _, t := range art.TagCloud() {
		tagged := articles.Tag(t.Tag)
		sm = append(sm, sitemap{
			Loc:        "http://" + r.Host + "/tag/" + t.Tag,
			Priority:   0.6 - float64(t.Wight)/10,
			LastMod:    lastMod(tagged[0].Date),
			ChangeFreq: "weekly",
		})
	}
	err := tmpl.ExecuteTemplate(w, "sitemap.tmpl", sm)
	if err != nil {
		log.Println(err)
	}
}
