package main

import (
	"fmt"
	"log"
	"net/http"
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
	for y, v := range blog.Enabled().Articles().YearMap() {
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
	perpage := conf.Blog.ArticlesPerPage
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
	e := blog.Enabled().Articles()
	p.URL = "http://" + r.Host
	p.Pager(r.URL.Query().Get("page"))
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

func (p *tagPage) Select(match []string) bool {
	s := match[0]
	tagged, ok := blog.Enabled().Articles().TagMap()[s]
	p.Articles = tagged
	p.Title = s
	return ok
}

type indexPage struct{ page }

func (p *indexPage) Select(_ []string) bool {
	p.Articles = blog.Enabled().Articles()
	return true
}

type slugPage struct{ page }

func (p *slugPage) Select(match []string) bool {
	slug := match[0]

	if a, ok := blog.Enabled()[slug]; ok {
		p.Title = a.Title
		p.Articles = append(p.Articles, a)
		p.Year = a.Year()
		p.Month = a.Month()
		return true
	}
	return false
}

type yearPage struct{ page }

func (p *yearPage) Select(match []string) bool {
	p.Year = atoiMust(match[0])
	p.Articles = blog.Enabled().Articles().Year(p.Year)
	p.Title = fmt.Sprint(p.Year)
	return true
}

type monthPage struct{ page }

func (p *monthPage) Select(match []string) bool {
	p.Year = atoiMust(match[0])
	p.Month = time.Month(atoiMust(match[1]))
	p.Articles = blog.Enabled().Articles().Year(p.Year).Month(p.Month)
	p.Title = fmt.Sprint(p.Month, p.Year)
	return true
}
