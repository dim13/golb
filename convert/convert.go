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

type Blog struct {
	Articles Articles
	Drafts   Articles
}

type Articles map[string]Article

type Article struct {
	Date     time.Time
	Title    string
	Slug     string
	Body     string
	Tags     Tags
	Author   string
	Comments Comments
	Moderate Comments
}

type Tags []string

type Comments []Comment

type Comment struct {
	Date    time.Time
	Name    string
	Email   string
	URL     string
	Comment string
}

var (
	input  string
	output string
)

func write(fname string, v interface{}) {
	w, err := os.Create(fname)
	if err != nil {
		log.Fatal("create ", err)
	}
	defer w.Close()
	enc := gob.NewEncoder(w)
	err = enc.Encode(v)
	if err != nil {
		log.Fatal("encode ", err)
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
		log.Fatal("parse time ", err)
	}
	return d
}

func (a Article) String() string {
	return fmt.Sprintf("%s %s %s", a.Date.Format(timeFormat), a.Title, a.Tags)
}

func (c Comment) String() string {
	return fmt.Sprintf("%s Commentar from %s", c.Date.Format(timeFormat), c.Name)
}

func getComments(db *sql.DB, id int) (C Comments, M Comments) {
	rows, err := db.Query("SELECT date,name,email,url,comment,enabled FROM comments WHERE article_id=?", id)
	if err != nil {
		log.Fatal("query comment ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			date    time.Time
			name    string
			email   []byte
			url     []byte
			comment string
			enabled bool
		)

		err := rows.Scan(&date, &name, &email, &url, &comment, &enabled)
		if err != nil {
			log.Fatal("scan comment ", err)
		}

		c := Comment{
			Date:    date,
			Name:    name,
			Email:   string(email),
			URL:     string(url),
			Comment: comment,
			Enabled: enabled,
		}

		fmt.Println(c)
		if enabled {
			C = append(C, c)
		} else {
			M = append(M, c)
		}
	}

	return C, M
}

func getArticles(db *sql.DB) (A Articles, D Articles) {
	rows, err := db.Query("SELECT id,date,title,uri,body,tags,enabled,author FROM articles")
	if err != nil {
		log.Fatal("query article ", err)
	}
	defer rows.Close()

	A = make(Articles)
	D = make(Articles)

	for rows.Next() {
		var (
			id      int
			date    time.Time
			title   string
			uri     string
			body    string
			tags    string
			enabled bool
			author  string
		)

		err := rows.Scan(&id, &date, &title, &uri, &body, &tags, &enabled, &author)
		if err != nil {
			log.Fatal("scan article ", err)
		}

		c, m := getComments(db, id)

		a := Article{
			Date:     date,
			Title:    title,
			Slug:     uri,
			Body:     body,
			Tags:     getTags(tags),
			Author:   author,
			Comments: c,
			Moderate: m,
		}

		fmt.Println(a)
		if enabled {
			A[uri] = a
		} else {
			D[uri] = a
		}
	}

	return A, D
}

func init() {
	flag.StringVar(&input, "input", "site.db", "input file (sqlite3)")
	flag.StringVar(&output, "output", "site.gob", "output file (gob)")
	flag.Parse()
}

func main() {
	db, err := sql.Open("sqlite3", input)
	if err != nil {
		log.Fatal("open ", err)
	}
	defer db.Close()
	a, d := getArticles(db)
	write(output, Blog{
		Articles: a,
		Drafts:   d,
	})
}
