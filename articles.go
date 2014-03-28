// Articles
package golb

import (
	"errors"
	"strings"
	"time"
)

type Articles []*Article

type Article struct {
	Date     time.Time
	Title    string
	Slug     string
	Body     string
	Tags     []string
	Enabled  bool
	Author   string
	Comments Comments
}

func (a *Article) makeSlug() {
	r := strings.NewReplacer(" ", "-")
	a.Slug = strings.ToLower(r.Replace(a.Title))
}

func (a *Article) Publish() {
	a.Date = time.Now()
	a.Enabled = true
}

func (a *Article) Suppress() {
	a.Enabled = false
}

func (a Articles) Len() int {
	return len(a)
}

func (a Articles) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Articles) Less(i, j int) bool {
	return a[i].Date.Before(a[j].Date)
}

func (a *Articles) Add(article *Article) {
	article.Date = time.Now()
	if article.Slug == "" {
		article.makeSlug()
	}
	*a = append(*a, article)
}

func (a Articles) Find(slug string) (*Article, error) {
	for i, _ := range a {
		if a[i].Slug == slug {
			return a[i], nil
		}
	}
	return &Article{}, errors.New("not found")
}

func (a Articles) CountTags() map[string]int {
	tags := make(map[string]int)
	for _, article := range a {
		for _, tag := range article.Tags {
			tags[tag]++
		}
	}
	return tags
}
