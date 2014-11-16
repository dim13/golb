package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/dim13/gold/articles"
	"github.com/dim13/gold/storage"
)

type page struct {
	Config    storage.Config
	URL       string
	Title     string
	Articles  articles.Articles
	Error     error
	PrevPage  int
	NextPage  int
	TagCloud  articles.TagCloud
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
	Articles articles.Articles
}

func (m byMonth) Len() int           { return len(m) }
func (m byMonth) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m byMonth) Less(i, j int) bool { return m[i].Month < m[j].Month }

func (y byYear) Len() int           { return len(y) }
func (y byYear) Swap(i, j int)      { y[i], y[j] = y[j], y[i] }
func (y byYear) Less(i, j int) bool { return y[i].Year < y[j].Year }

func atoiMust(s string) int {
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
	for y, v := range art.Enabled().YearMap() {
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

func getPage(u *url.URL) int {
	if page, ok := u.Query()["page"]; ok {
		return atoiMust(page[0])
	}
	return 1
}

func (p *page) Pager(pg, pp int) {
	if pg <= 1 {
		pg = 1
	} else {
		p.PrevPage = pg - 1
	}

	n := len(p.Articles)
	lastpage := n/pp + 1

	if pg >= lastpage {
		pg = lastpage
	} else {
		p.NextPage = pg + 1
	}

	from := (pg - 1) * pp
	to := from + pp - 1

	if to > n {
		to = n
	}

	p.Articles = p.Articles[from:to]
}

func (p page) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e := art.Enabled()
	p.URL = "http://" + r.Host
	p.Pager(getPage(r.URL), conf.Blog.ArticlesPerPage)
	p.TagCloud = e.TagCloud()
	p.Config = conf
	p.MakeArchive()
	p.FirstYear = e.Tail().Year()
	p.LastYear = e.Head().Year()

	err := tmpl.ExecuteTemplate(w, "index.tmpl", p)
	if err != nil {
		log.Println(err)
	}
}

type tagPage struct{ page }

func (p *tagPage) Select(match []string) {
	s := match[0]
	p.Articles = art.Tag(s)
	p.Title = s
}

type indexPage struct{ page }

func (p *indexPage) Select(_ []string) {
	p.Articles = art.Enabled()
}

type slugPage struct{ page }

func (p *slugPage) Select(match []string) {
	if a, ok := art.Enabled().Find(match[0]); ok {
		p.Title = a.Title
		p.Articles = articles.Articles{a}
		p.Year = a.Year()
		p.Month = a.Month()
	}
}

type yearPage struct{ page }

func (p *yearPage) Select(match []string) {
	p.Year = atoiMust(match[0])
	p.Articles = art.Enabled().Year(p.Year)
	p.Title = fmt.Sprint(p.Year)
}

type monthPage struct{ page }

func (p *monthPage) Select(match []string) {
	p.Year = atoiMust(match[0])
	p.Month = time.Month(atoiMust(match[1]))
	p.Articles = art.Enabled().Year(p.Year).Month(p.Month)
	p.Title = fmt.Sprint(p.Month, p.Year)
}
