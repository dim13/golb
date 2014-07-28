package articles

import (
	"fmt"
	"strings"
	"time"

	"github.com/dim13/gold/storage"
)

const (
	TimeFormat = "January 2, 2006"
	readMore   = "<!--readmore-->"
)

var storageFile string

type Articles []*Article
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

type YearMap map[int]Articles
type MonthMap map[int]Articles

func MakeSlug(title string) string {
	r := strings.NewReplacer(" ", "-")
	return r.Replace(strings.TrimSpace(title))
}

func MakeTitle(slug string) string {
	r := strings.NewReplacer("-", " ")
	return r.Replace(strings.TrimSpace(slug))
}

func SetStorage(file string) {
	storageFile = file
}

func (a *Articles) Load() error {
	return storage.Load(storageFile, a)
}

func (a *Articles) Store() error {
	return storage.Store(storageFile, a)
}

func (a *Article) Publish() {
	a.Date = time.Now()
	a.Enabled = true
}

func (a *Article) Suppress() {
	a.Enabled = false
}

func (a Articles) Len() int           { return len(a) }
func (a Articles) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Articles) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }

func (a *Article) ChckSlug() {
	if a.Slug == "" {
		a.Slug = MakeSlug(a.Title)
	}
}

func (a *Articles) Add(art Article) {
	art.ChckSlug()
	if ar, _, ok := a.Find(art.Slug); ok {
		/* flound slug, update */
		art.Date = ar.Date
		*ar = art
	} else {
		/* no slug, add new */
		art.Date = time.Now()
		*a = append(Articles{&art}, *a...)
	}
	a.Store()
}

func (a *Articles) Delete(art Article) {
	if _, i, ok := a.Find(art.Slug); ok {
		(*a)[i] = nil
		*a = append((*a)[:i], (*a)[i+1:]...)
	}
}

func (a Articles) Find(slug string) (*Article, int, bool) {
	for i, ar := range a {
		if ar.Slug == slug {
			return ar, i, true
		}
	}
	return nil, 0, false
}

// Format Date with TimeFormat
func (a Article) PostDate() string {
	return a.Date.Local().Format(TimeFormat)
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
