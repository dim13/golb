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
	if err := d.Read(); err != nil {
		log.Println(err)
	}

	a := &golb.Article{
		Title: "Test title",
		Body: "empty body",
		Tags: []string{"no", "tags", "at all"},
		Author: "me@example.com",
	}
	err := d.Articles.Add(a)
	if err != nil {
		log.Println(err)
	}

	c := &golb.Comment{
		Name: "anonymous coward",
		Email: "none@example.com",
		URL: "http://example.com",
		Comment: "empty",
	}
	a.Comments.Add(c)
	c.Publish()

	z, err := d.Articles.Find("test-title")
	if err == nil {
		z.Publish()
	}

	conf, err := golb.ReadConf("config.ini")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(conf)

	if err := d.Write(); err != nil {
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
