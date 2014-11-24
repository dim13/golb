package blog

import (
	"github.com/dim13/gold/storage"
)

var storageFile string

type Blog map[string]Article

func SetStorage(file string) {
	storageFile = file
}

func Load() (Blog, error) {
	b := make(Blog)
	err := storage.Load(storageFile, &b)
	return b, err
}

func (b Blog) Store() error {
	return storage.Store(storageFile, b)
}

func (b Blog) Add(a Article) {
	slug := a.Slug()

	if ar, ok := b[slug]; ok {
		/* found, preserve date */
		a.Date = ar.Date
	}

	b[slug] = a
}

func (b Blog) Delete(a Article) {
	slug := a.Slug()

	if _, ok := b[slug]; ok {
		delete(b, slug)
	}
}

func (b Blog) Find(slug string) (Article, bool) {
	a, ok := b[slug]
	return a, ok
}

func (b Blog) Articles() (a Articles) {
	for _, v := range b {
		a = append(a, v)
	}
	return a.Sort()
}

func (b Blog) Enabled() Blog {
	e := make(Blog)
	for k, v := range b {
		if v.Enabled {
			e[k] = v
		}
	}
	return e
}
