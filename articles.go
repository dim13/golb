// JSON Data Storage
package golb

import (
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

func (a *Article) Add(c *Comment) {
	c.Date = time.Now()
	a.Comments = append(a.Comments, c)
}

func (a *Article) Enable() {
	a.Date = time.Now()
	a.Enabled = true
}

func (a *Article) Disable() {
	a.Enabled = false
}

func (a Articles) Len() int {
	return len(a)
}

func (a Articles) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Articles) Less(i, j int) bool {
	return a[i].Date.After(a[j].Date)
}
