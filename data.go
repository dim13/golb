// JSON Data Storage
package golb

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"sort"
	"time"
)

type Data struct {
	Articles Articles
	fileName     string
}

func Open(name string) *Data {
	d := new(Data)
	d.fileName = name
	return d
}

func (d *Data) Read() error {
	data, err := ioutil.ReadFile(d.fileName)
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
	return ioutil.WriteFile(d.fileName, data, 0644)
}

func (d *Data) Find(slug string) (*Article, error) {
	for _, a := range d.Articles {
		if a.Slug == slug {
			return a, nil
		}
	}
	return &Article{}, errors.New("not found")
}

func (d *Data) Add(a *Article) {
	a.Date = time.Now()
	if a.Slug == "" {
		a.makeSlug()
	}
	d.Articles = append(d.Articles, a)
}

func (d *Data) CountTags() map[string]int {
	tags := make(map[string]int)
	for _, a := range d.Articles {
		for _, t := range a.Tags {
			tags[t]++
		}
	}
	return tags
}
