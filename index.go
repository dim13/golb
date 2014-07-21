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

type Page struct {
	Config    storage.Config
	Title     string
	Articles  articles.Articles
	Error     error
	PrevPage  int
	NextPage  int
	TagCloud  articles.TagCloud
	Year      int
	Month     time.Month
	Archive   []Archive
	FirstYear int
	LastYear  int
}

type ByYear []Archive
type Archive struct {
	Year  int
	Count int
	Month []Month
}

type ByMonth []Month
type Month struct {
	Month    time.Month
	Year     int
	Count    int
	Articles articles.Articles
}

func (m ByMonth) Len() int           { return len(m) }
func (m ByMonth) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m ByMonth) Less(i, j int) bool { return m[i].Month < m[j].Month }

func (y ByYear) Len() int           { return len(y) }
func (y ByYear) Swap(i, j int)      { y[i], y[j] = y[j], y[i] }
func (y ByYear) Less(i, j int) bool { return y[i].Year < y[j].Year }

func atoiMust(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func getPage(u url.URL) int {
	if page, ok := u.Query()["page"]; ok {
		if pg, err := strconv.Atoi(page[0]); err == nil {
			return pg
		}
	}
	return 1
}

func (p *Page) MakeArchive() {
	if p.Year == 0 {
		p.Year = p.Articles.Head().Year()
	}
	if p.Month == 0 {
		p.Month = p.Articles.Head().Month()
	}
	for y, v := range art.Enabled().YearMap() {
		year := Archive{
			Year:  y,
			Count: len(v),
		}
		if p.Year == y {
			for m, v := range v.MonthMap() {
				month := Month{
					Year:  y,
					Month: time.Month(m),
					Count: len(v),
				}
				if p.Month == time.Month(m) {
					month.Articles = v
				}
				year.Month = append(year.Month, month)
			}
			sort.Sort(ByMonth(year.Month))
		}
		p.Archive = append(p.Archive, year)
	}
	sort.Sort(sort.Reverse(ByYear(p.Archive)))
}

func (p *Page) Pager(pg, pp int) {
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

func (p Page) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.Pager(getPage(*r.URL), conf.Blog.ArticlesPerPage)
	p.TagCloud = art.TagCloud()
	p.Config = conf
	p.MakeArchive()
	p.FirstYear = art.Enabled().Tail().Year()
	p.LastYear = art.Enabled().Head().Year()

	err := tmpl.ExecuteTemplate(w, "index.tmpl", p)
	if err != nil {
		log.Println(err)
	}
}

type TagPage struct{ Page }

func (p *TagPage) Select(match []string) {
	s := match[0]
	p.Articles = art.Tag(s)
	p.Title = s
}

type IndexPage struct{ Page }

func (p *IndexPage) Select(_ []string) {
	p.Articles = art.Enabled()
}

type SlugPage struct{ Page }

func (p *SlugPage) Select(match []string) {
	if a, ok := art.Find(match[0]); ok {
		p.Title = a.Title
		p.Articles = articles.Articles{a}
		p.Year = a.Year()
		p.Month = a.Month()
	}
}

type YearPage struct{ Page }

func (p *YearPage) Select(match []string) {
	p.Year = atoiMust(match[0])
	p.Articles = art.Year(p.Year)
	p.Title = fmt.Sprint(p.Year)
}

type MonthPage struct{ Page }

func (p *MonthPage) Select(match []string) {
	p.Year = atoiMust(match[0])
	p.Month = time.Month(atoiMust(match[1]))
	p.Articles = art.Year(p.Year).Month(p.Month)
	p.Title = fmt.Sprint(p.Month, p.Year)
}
