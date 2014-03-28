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
