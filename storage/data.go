package storage

import (
	"github.com/dim13/gold/articles"
	"sort"
)

type Data struct {
	Articles articles.Articles
	fileName string
	TagMap   articles.TagMap
	Tags     articles.TagCount
}

func Open(name string) *Data {
	d := new(Data)
	d.fileName = name
	return d
}

func (d *Data) Read() error {
	return ReadJson(d.fileName, &d.Articles)
}

func (d *Data) Write() error {
	sort.Sort(d.Articles)
	for i, _ := range d.Articles {
		sort.Sort(d.Articles[i].Comments)
	}
	return WriteJson(d.fileName, &d.Articles)
}

func (d *Data) AddArticle(a *articles.Article) {
	d.Articles.Add(a)
	d.Tags = d.Articles.CountTags()
	d.TagMap = d.Articles.TagMap()
	d.Write()
}
