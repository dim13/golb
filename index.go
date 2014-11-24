package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/dim13/gold/blog"
	"github.com/dim13/gold/storage"
)

type page struct {
	Config    storage.Config
	URL       string
	Title     string
	Articles  blog.Articles
	Error     error
	PrevPage  int
	NextPage  int
	TagCloud  blog.TagCloud
	Year      int
	Month     time.Month
	Archive   []year
	FirstYear int
	LastYear  int
}

type byYear []year
type year struct {
	Year  int
	Count int
	Month []month
}

type byMonth []month
type month struct {
	Month    time.Month
	Year     int
	Count    int
	Articles blog.Articles
}

func (m byMonth) Len() int           { return len(m) }
func (m byMonth) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m byMonth) Less(i, j int) bool { return m[i].Month < m[j].Month }

func (y byYear) Len() int           { return len(y) }
func (y byYear) Swap(i, j int)      { y[i], y[j] = y[j], y[i] }
func (y byYear) Less(i, j int) bool { return y[i].Year < y[j].Year }

func atoiMust(s string) int {
	if s == "" {
		return 0
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Println(err)
		return 0
	}
	return i
}

func (p *page) MakeArchive() {
	if p.Year == 0 {
		p.Year = p.Articles.Head().Year()
	}
	if p.Month == 0 {
		p.Month = p.Articles.Head().Month()
	}
	for y, v := range Blog.Enabled().Articles().YearMap() {
		year := year{
			Year:  y,
			Count: len(v),
		}
		if p.Year == y {
			for m, v := range v.MonthMap() {
				month := month{
					Year:  y,
					Month: time.Month(m),
					Count: len(v),
				}
				if p.Month == time.Month(m) {
					month.Articles = v
				}
				year.Month = append(year.Month, month)
			}
			sort.Sort(byMonth(year.Month))
		}
		p.Archive = append(p.Archive, year)
	}
	sort.Sort(sort.Reverse(byYear(p.Archive)))
}

func (p *page) Pager(pg string) {
	perpage := Conf.Blog.ArticlesPerPage
	count := len(p.Articles)
	last := count/perpage + 1
	curr := atoiMust(pg)

	if curr <= 1 {
		curr = 1
	} else {
		p.PrevPage = curr - 1
	}

	if curr >= last {
		curr = last
	} else {
		p.NextPage = curr + 1
	}

	from := (curr - 1) * perpage
	to := from + perpage - 1

	if to > count {
		to = count
	}

	p.Articles = p.Articles[from:to]
}

func (p page) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.URL = "http://" + r.Host
	p.Pager(r.URL.Query().Get("page"))
	p.TagCloud = Blog.Enabled().TagCloud()
	p.Config = Conf
	p.MakeArchive()

	a := Blog.Enabled().Articles()
	p.FirstYear = a.Tail().Year()
	p.LastYear = a.Head().Year()

	err := tmpl.ExecuteTemplate(w, "index.tmpl", p)
	if err != nil {
		log.Println(err)
	}
}

func tagHandler(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get(":tag")
	tagged, _ := Blog.Enabled().TagMap()[tag]
	pg := page{
		Articles: tagged,
		Title:    tag,
	}
	pg.ServeHTTP(w, r)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	pg := page{
		Articles: Blog.Enabled().Articles(),
	}
	pg.ServeHTTP(w, r)
}

func slugHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get(":slug")
	a, _ := Blog.Enabled()[slug]
	pg := page{
		Title:    a.Title,
		Articles: blog.Articles{a},
		Year:     a.Year(),
		Month:    a.Month(),
	}
	pg.ServeHTTP(w, r)
}

func yearHandler(w http.ResponseWriter, r *http.Request) {
	year := atoiMust(r.URL.Query().Get(":year"))
	pg := page{
		Year:     year,
		Articles: Blog.Enabled().Articles().Year(year),
		Title:    fmt.Sprint(year),
	}
	pg.ServeHTTP(w, r)
}

func monthHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	year := atoiMust(q.Get(":year"))
	month := time.Month(atoiMust(q.Get(":month")))
	pg := page{
		Year:     year,
		Month:    month,
		Articles: Blog.Enabled().Articles().Year(year).Month(month),
		Title:    fmt.Sprint(month, year),
	}
	pg.ServeHTTP(w, r)
}
