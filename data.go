// JSON Data Storage
package golb

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"
)

type Comment struct {
	Date    time.Time
	Name    string
	Email   string
	URL     string
	Comment string
	Enabled bool
}

type Article struct {
	Date     time.Time
	Title    string
	Slug     string
	Body     string
	Tags     []string
	Enabled  bool
	Author   string
	Comments []Comment
}

type Data struct {
	Articles []Article
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
	data, err := json.MarshalIndent(d.Articles, "", "	")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(d.Name, data, 0644)
}

func (d *Data) FindArticle(slug string) (Article, error) {
	for _, a := range d.Articles {
		if a.Slug == slug {
			return a, nil
		}
	}
	return Article{}, errors.New("not found")
}
