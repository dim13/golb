// JSON Data Storage
package golb

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"sort"
	"strings"
	"time"
)

type Comments []*Comment

type Comment struct {
	Date    time.Time
	Name    string
	Email   string
	URL     string
	Comment string
	Enabled bool
}

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

type Data struct {
	Articles Articles
	Name     string
}

func Open(name string) *Data {
	d := new(Data)
	d.Name = name
	return d
}

func (d *Data) Read() error {
	data, err := ioutil.ReadFile(d.Name)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &d.Articles)
}

func (d *Data) Write() error {
	sort.Sort(d.Articles)
	for i, _ := range d.Articles {
		sort.Sort(d.Articles[i].Comments)
	}
	data, err := json.MarshalIndent(d.Articles, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(d.Name, data, 0644)
}

func (d *Data) FindArticle(slug string) (*Article, error) {
	for _, a := range d.Articles {
		if a.Slug == slug {
			return a, nil
		}
	}
	return &Article{}, errors.New("not found")
}

func (a *Article) MakeSlug() {
	r := strings.NewReplacer(" ", "-")
	a.Slug = strings.ToLower(r.Replace(a.Title))
}

func (d *Data) AddArticle(a *Article) {
	a.Date = time.Now()
	a.MakeSlug()
	d.Articles = append(d.Articles, a)
}

func (a *Article) AddComment(c *Comment) {
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

func (c *Comment) Enable() {
	c.Enabled = true
}

func (c *Comment) Disable() {
	c.Enabled = false
}

func (c Comments) Len() int {
	return len(c)
}

func (c Comments) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c Comments) Less(i, j int) bool {
	return c[i].Date.After(c[j].Date)
}
