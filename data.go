package gold

import (
	"github.com/dim13/jsondb"
	"sort"
)

type Data struct {
	Articles Articles
	fileName string
	TagMap   TagMap
	Tags     TagCount
}

func Open(name string) *Data {
	d := new(Data)
	d.fileName = name
	return d
}

func (d *Data) Read() error {
	return jsondb.Read(d.fileName, &d.Articles)
}

func (d *Data) Write() error {
	sort.Sort(d.Articles)
	for i, _ := range d.Articles {
		sort.Sort(d.Articles[i].Comments)
	}
	return jsondb.Write(d.fileName, &d.Articles)
}

func (d *Data) AddArticle(a *Article) {
	d.Articles.Add(a)
	d.Tags = d.Articles.CountTags()
	d.TagMap = d.Articles.TagMap()
	d.Write()
}
