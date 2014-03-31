package golb

import (
	"encoding/json"
	"io/ioutil"
	"sort"
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

func (d *Data) Add(a *Article) {
	d.Articles.Add(a)
	d.Tags = d.Articles.CountTags()
	d.TagMap = d.Articles.TagMap()
	d.Write()
}
