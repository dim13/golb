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

func Load(fname string) (*Data, error) {
	d := &Data{fileName: fname}
	return d, load(d.fileName, &d.Articles)
}

func (d *Data) Store() error {
	sort.Sort(sort.Reverse(d.Articles))
	for i, _ := range d.Articles {
		sort.Sort(d.Articles[i].Comments)
	}
	return store(d.fileName, &d.Articles)
}

func (d *Data) AddArticle(a *articles.Article) {
	d.Articles.Add(a)
	d.Tags = d.Articles.CountTags()
	d.TagMap = d.Articles.TagMap()
	d.Store()
}
