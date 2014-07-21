package articles

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	TimeFormat = "January 2, 2006"
	readMore   = "<!--readmore-->"
)

type Articles []*Article
type Article struct {
	Date     time.Time
	Edit     time.Time
	Title    string
	Slug     string
	Body     string
	Tags     Tags
	Enabled  bool
	Author   string
	Comments Comments
}

type YearMap map[int]Articles
type MonthMap map[int]Articles

func (a *Article) makeSlug() {
	r := strings.NewReplacer(" ", "-")
	a.Slug = r.Replace(strings.TrimSpace(a.Title))
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

func (a Articles) Len() int           { return len(a) }
func (a Articles) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Articles) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }

func (a *Articles) Add(article *Article) error {
	article.Date = time.Now()
	if article.Slug == "" {
		article.makeSlug()
	}
	if _, ok := a.Find(article.Slug); ok {
		return errors.New("duplicate slug " + article.Slug)
	}
	*a = append(*a, article)
	return nil
}

func (a Articles) Update(article *Article) error {
	article.Edit = time.Now()
	if article.Slug == "" {
		article.makeSlug()
	}
	for i, ar := range a {
		if ar.Slug == article.Slug {
			article.Date = ar.Date
			a[i] = ar
			return nil
		}
	}
	return a.Add(article)
}

func (a Articles) Find(slug string) (*Article, bool) {
	for _, ar := range a {
		if ar.Slug == slug {
			return ar, true
		}
	}
	return nil, false
}

// Format Date with TimeFormat
func (a Article) PostDate() string {
	return a.Date.Local().Format(TimeFormat)
}

func (a Article) EditDate() string {
	return a.Edit.Local().Format(TimeFormat)
}

func (a Article) Edited() bool {
	return !a.Edit.IsZero()
}

// Format Date for RSS
func (a Article) RssDate() string {
	return a.Date.Local().Format(time.RFC1123Z)
}

func (a Article) Spoiler() string {
	if i := strings.Index(a.Body, readMore); i > 0 {
		return a.Body[:i]
	}
	return a.Body
}

func (a Article) HasMore() bool {
	return strings.Contains(a.Body, readMore)
}

func (a Article) Year() int {
	return a.Date.Year()
}

func (a Article) Month() time.Month {
	return a.Date.Month()
}

func (a Articles) Year(year int) (A Articles) {
	if year == 0 {
		year = time.Now().Year()
	}
	for _, v := range a {
		if v.Date.Year() == year {
			A = append(A, v)
		}
	}
	return
}

func (a Articles) Month(month time.Month) (A Articles) {
	if month == 0 {
		month = time.Now().Month()
	}
	for _, v := range a {
		if v.Date.Month() == month {
			A = append(A, v)
		}
	}
	return
}

func (a Articles) Enabled() (A Articles) {
	for _, v := range a {
		if v.Enabled {
			A = append(A, v)
		}
	}
	return
}

func (a Articles) Skip(n int) Articles {
	if n > len(a) {
		return nil
	}
	return a[n:]
}

func (a Articles) Limit(n int) Articles {
	if n > len(a) {
		n = len(a)
	}
	return a[:n]
}

func (a Articles) Head() Article {
	if len(a) > 0 {
		return *a[0]
	}
	return Article{}
}

func (a Articles) Tail() Article {
	if len(a) > 0 {
		return *a[len(a)-1]
	}
	return Article{}
}

func (a Articles) YearMap() YearMap {
	ym := make(YearMap)
	for _, v := range a {
		y := v.Date.Year()
		ym[y] = append(ym[y], v)
	}
	return ym
}

func (a Articles) MonthMap() MonthMap {
	mm := make(MonthMap)
	for _, v := range a {
		m := int(v.Date.Month())
		mm[m] = append(mm[m], v)
	}
	return mm
}

func (a Article) FullPath() string {
	return fmt.Sprintf("/%.4d/%.2d/%s",
		a.Date.Year(), a.Date.Month(), a.Slug)
}
