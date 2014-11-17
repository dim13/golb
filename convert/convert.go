// Convert BlogSum DB (sqlite3) to GOB
package main

import (
	"database/sql"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const timeFormat = "2006-Jan-02"

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
	input  string
	output string
)

func writeGob(fname string, v interface{}) {
	w, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()
	enc := gob.NewEncoder(w)
	err = enc.Encode(v)
	if err != nil {
		log.Fatal(err)
	}
}

func getTags(tags string) Tags {
	t := strings.Split(tags, ",")
	for i := range t {
		t[i] = strings.TrimSpace(t[i])
	}
	return t
}

func getDate(date string) time.Time {
	d, err := time.Parse("2006-01-02 15:04:05", date)
	if err != nil {
		log.Fatal(err)
	}
	return d
}

func (a Article) String() string {
	return fmt.Sprintf("%s %s %s", a.Date.Format(timeFormat), a.Slug, a.Tags)
}

func (c Comment) String() string {
	return fmt.Sprintf("%s Commentar from %s", c.Date.Format(timeFormat), c.Name)
}

func getComments(db *sql.DB, id int) (C Comments) {
	rows, err := db.Query("SELECT date,name,email,url,comment,enabled FROM comments WHERE article_id=?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			date    string
			name    string
			email   []byte
			url     []byte
			comment string
			enabled bool
		)

		err := rows.Scan(&date, &name, &email, &url, &comment, &enabled)
		if err != nil {
			log.Fatal(err)
		}

		c := Comment{
			Date:    getDate(date),
			Name:    name,
			Email:   string(email),
			URL:     string(url),
			Comment: comment,
			Enabled: enabled,
		}

		fmt.Println(c)
		C = append(C, c)
	}

	return C
}

func getArticles(db *sql.DB) (A Articles) {
	rows, err := db.Query("SELECT id,date,title,uri,body,tags,enabled,author FROM articles")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id      int
			date    string
			title   string
			uri     string
			body    string
			tags    string
			enabled bool
			author  string
		)

		err := rows.Scan(&id, &date, &title, &uri, &body, &tags, &enabled, &author)
		if err != nil {
			log.Fatal(err)
		}

		a := Article{
			Date:     getDate(date),
			Title:    title,
			Slug:     uri,
			Body:     body,
			Tags:     getTags(tags),
			Enabled:  enabled,
			Author:   author,
			Comments: getComments(db, id),
		}

		fmt.Println(a)
		A = append(A, a)
	}

	return A
}

func init() {
	flag.StringVar(&input, "input", "site.db", "input file (sqlite3)")
	flag.StringVar(&output, "output", "site.gob", "output file (gob)")
	flag.Parse()
}

func main() {
	db, err := sql.Open("sqlite3", input)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	writeGob(output, getArticles(db))
}
