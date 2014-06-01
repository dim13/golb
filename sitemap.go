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

func (sm *SiteMap) Select(match []string) {
	*sm = append(*sm, Sitemap{
		Loc:      conf.Blog.Url,
		Priority: 1.0,
	})
	for _, a := range data.Articles.Enabled() {
		*sm = append(*sm, Sitemap{
			Loc:        conf.Blog.Url + "/" + a.Slug,
			Priority:   0.8,
			Date:       a.Date,
			ChangeFreq: "monthly",
		})
	}
	for _, t := range data.Articles.TagCloud() {
		*sm = append(*sm, Sitemap{
			Loc:      conf.Blog.Url + "/tag/" + t.Tag,
			Priority: 0.6 - float64(t.Wight)/10,
		})
	}
}

func (sm *SiteMap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "sitemap.tmpl", sm)
	if err != nil {
		log.Fatal(err)
	}
}

func (sm *SiteMap) Store(r *http.Request) {}
