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

	"github.com/dim13/gold/blog"

	_ "github.com/mattn/go-sqlite3"
)

const timeFormat = "2006-Jan-02"

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

func getTags(tags string) blog.Tags {
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

func getComments(db *sql.DB, id int) (C blog.Comments, M blog.Comments) {
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

		c := blog.Comment{
			Date:    date,
			Name:    name,
			Email:   string(email),
			URL:     string(url),
			Comment: comment,
		}

		fmt.Printf("%s Commentar from %s\n",
			c.Date.Format(timeFormat), c.Name)

		if enabled {
			C = append(C, c)
		} else {
			M = append(M, c)
		}
	}

	return C, M
}

func getArticles(db *sql.DB) (A blog.Items, D blog.Items) {
	rows, err := db.Query("SELECT id,date,title,uri,body,tags,enabled,author FROM articles")
	if err != nil {
		log.Fatal("query article ", err)
	}
	defer rows.Close()

	A = make(blog.Items)
	D = make(blog.Items)

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

		a := blog.Article{
			Date:     date,
			Title:    title,
			Slug:     uri,
			Body:     body,
			Tags:     getTags(tags),
			Author:   author,
			Comments: c,
			Moderate: m,
		}

		fmt.Printf("%s %s %s\n",
			a.Date.Format(timeFormat), a.Title, a.Tags)

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
	write(output, blog.Blog{
		Public: a,
		Draft:  d,
	})
}
