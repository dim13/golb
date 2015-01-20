package blog

import (
	"github.com/dim13/gold/storage"
)

var storageFile string

type Blog struct {
	Public map[string]Article
	Draft  map[string]Article
}

//type Blog map[string]Article

func SetStorage(file string) {
	storageFile = file
}

func Load() (Blog, error) {
	var b Blog
	err := storage.Load(storageFile, &b)
	return b, err
}

func (b Blog) Store() error {
	return storage.Store(storageFile, b)
}

func (b Blog) Add(a Article) {
	if ar, ok := b.Draft[a.Slug]; ok {
		/* found, preserve date */
		a.Date = ar.Date
	}
	b.Draft[a.Slug] = a
}

func (b Blog) Delete(a Article) {
	delete(b.Draft, a.Slug)
}

func (b Blog) Find(slug string) (Article, bool) {
	a, ok := b.Public[slug]
	return a, ok
}

func (b Blog) Drafts() (a Articles) {
	for _, v := range b.Draft {
		a = append(a, v)
	}
	return a.Sort()
}

func (b Blog) Articles() (a Articles) {
	for _, v := range b.Public {
		a = append(a, v)
	}
	return a.Sort()
}

func (b Blog) Publish(slug string) {
	b.Public[slug] = b.Draft[slug]
	//b.Public[slug].Date = time.Now()
	delete(b.Draft, slug)
}

func (b Blog) Concial(slug string) {
	b.Draft[slug] = b.Public[slug]
	delete(b.Public, slug)
}
