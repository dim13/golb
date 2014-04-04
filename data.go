package gold

import (
	"sort"
	"github.com/dim13/jsondb"
)

type Data struct {
	Articles Articles
	fileName string
	TagMap   TagMap
	Tags     Tags
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
