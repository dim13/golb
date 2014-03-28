// Test APP
package main

import (
	"fmt"
	"github.com/dim13/golb"
	"log"
	"net/http"
	"time"
)

const listen = ":8000"

func root(w http.ResponseWriter, r *http.Request, s []string) {
	fmt.Fprint(w, s)
}

func main() {
	d := golb.Open("test.json")
	c := golb.Comment{
		Date: time.Now(),
		Name: "anonymous coward",
		Email: "none@example.com",
		URL: "http://example.com",
		Comment: "empty",
		Enabled: true,
	}
	a := golb.Article{
		Date: time.Now(),
		Title: "Test title",
		Slug: "test-title",
		Body: "empty body",
		Tags: []string{"no", "tags", "at all"},
		Enabled: true,
		Author: "me@example.com",
	}
	a.Comments = append(a.Comments, c)
	d.Articles = append(d.Articles, a)

	if err := d.Write(); err != nil {
		log.Fatal(err)
	}
	if err := d.Read(); err != nil {
		log.Fatal(err)
	}

	re := new(golb.ReHandler)
	re.AddRoute("^/(\\d+)/(.*)$", root)
	if err := http.ListenAndServe(listen, re); err != nil {
		log.Fatal(err)
	}
}
