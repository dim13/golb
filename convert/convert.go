// Convert BlogSum DB (sqlite3) to JSON
package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

/*
BlogSum DB Schema:

CREATE TABLE articles (
	id integer primary key,
	date date,
	title text,
	uri text,
	body text,
	tags text,
	enabled boolean,
	author text);

CREATE TABLE comments (
	id integer primary key,
	article_id integer,
	date date,
	name text,
	email text,
	url text,
	comment text,
	enabled boolean);
*/

type Articles []Article
type Article struct {
	Date     time.Time
	Title    string
	Slug     string // Uri
	Body     string
	Tags     Tags
	Enabled  bool
	Author   string
	Comments Comments
}
type Tags []string
type Comments []Comment
type Comment struct {
	Date    time.Time
	Name    string
	Email   string
	URL     string
	Comment string
	Enabled bool
}

var (
	format = "2006-01-02 15:04:05"
	input  = "site.db"
	output = "site.json"
)

func main() {
	var A Articles

	db, err := sql.Open("sqlite3", input)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id,date,title,uri,body,tags,enabled,author FROM articles")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var date string
		var title string
		var uri string
		var body string
		var tags string
		var enabled bool
		var author string

		var C Comments

		err := rows.Scan(&id, &date, &title, &uri, &body, &tags, &enabled, &author)
		if err != nil {
			log.Fatal(err)
		}

		t := strings.Split(tags, ",")
		for i := range t {
			t[i] = strings.TrimSpace(t[i])
		}
		log.Println(uri, t)

		d, err := time.Parse(format, date)
		if err != nil {
			log.Fatal(err)
		}

		crows, err := db.Query("SELECT date,name,email,url,comment,enabled FROM comments WHERE article_id=?", id)
		if err != nil {
			log.Fatal(err)
		}
		defer crows.Close()

		for crows.Next() {
			var date string
			var name string
			var email []byte
			var url []byte
			var comment string
			var enabled bool

			err := crows.Scan(&date, &name, &email, &url, &comment, &enabled)
			if err != nil {
				log.Fatal(err)
			}

			d, err := time.Parse(format, date)
			if err != nil {
				log.Fatal(err)
			}

			c := Comment{
				Date:    d,
				Name:    name,
				Email:   string(email),
				URL:     string(url),
				Comment: comment,
				Enabled: enabled,
			}

			C = append(C, c)
		}

		a := Article{
			Date:     d,
			Title:    title,
			Slug:     uri,
			Body:     body,
			Tags:     t,
			Enabled:  enabled,
			Author:   author,
			Comments: C,
		}

		A = append(A, a)
	}

	data, err := json.MarshalIndent(A, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(output, data, 0664)
}
