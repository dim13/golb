// Test APP
package main

import (
	"fmt"
	"github.com/dim13/golb"
	"log"
	"net/http"
)

const listen = ":8000"

func root(w http.ResponseWriter, r *http.Request, s []string) {
	fmt.Fprint(w, s)
}

func main() {
	d := golb.Open("test.json")

	a := &golb.Article{
		Title: "Test title",
		Body: "empty body",
		Tags: []string{"no", "tags", "at all"},
		Author: "me@example.com",
	}
	d.Articles = d.Articles.Add(a)
	a.Publish()

	c := &golb.Comment{
		Name: "anonymous coward",
		Email: "none@example.com",
		URL: "http://example.com",
		Comment: "empty",
	}
	a.Comments = a.Comments.Add(c)
	c.Publish()

	if err := d.Write(); err != nil {
		log.Fatal(err)
	}
	if err := d.Read(); err != nil {
		log.Fatal(err)
	}

	/*
	re := new(golb.ReHandler)
	re.AddRoute("^/(\\d+)/(.*)$", root)
	if err := http.ListenAndServe(listen, re); err != nil {
		log.Fatal(err)
	}
	 */
}
