// Test APP
package main

import (
	"fmt"
	"github.com/dim13/gold"
	"log"
	"net/http"
)

const listen = ":8000"

func root(w http.ResponseWriter, r *http.Request, s []string) {
	fmt.Fprint(w, s)
}

func main() {
	d := gold.Open("test.json")
	if err := d.Read(); err != nil {
		log.Println(err)
	}

	a := &gold.Article{
		Title: "Test title",
		Body: "empty body",
		Tags: []string{"no", "tags", "at all"},
		Author: "me@example.com",
	}
	err := d.Articles.Add(a)
	if err != nil {
		log.Println(err)
	}

	c := &gold.Comment{
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

	conf, err := gold.ReadConf("config.ini")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(conf)

	if err := d.Write(); err != nil {
		log.Fatal(err)
	}

	/*
	re := new(gold.ReHandler)
	re.AddRoute("^/(\\d+)/(.*)$", root)
	if err := http.ListenAndServe(listen, re); err != nil {
		log.Fatal(err)
	}
	 */
}
