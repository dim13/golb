package gold

import (
	"errors"
	"strings"
	"time"
)

type Articles []*Article
type Tags []string

type Article struct {
	Date     time.Time
	Title    string
	Slug     string
	Body     string
	Tags     Tags
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

func (a *Article) AddComment(c *Comment) {
	a.Comments.Add(c)
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

func (a *Articles) Add(article *Article) error {
	article.Date = time.Now()
	if article.Slug == "" {
		article.makeSlug()
	}
	_, err := a.Find(article.Slug)
	if err == nil {
		return errors.New("duplicate slug " + article.Slug)
	}
	*a = append(*a, article)
	return nil
}

func (a Articles) Find(slug string) (*Article, error) {
	for i, _ := range a {
		if a[i].Slug == slug {
			return a[i], nil
		}
	}
	return &Article{}, errors.New("not found")
}

func (a Articles) Page(page, base int) Articles {
	return a[page*base : page*base+base]
}

func (t Tags) String() string {
	return strings.Join(t, ",")
}
