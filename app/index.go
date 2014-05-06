package main

import (
	"fmt"
	"github.com/dim13/gold"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

type Page struct {
	Config   *gold.Config
	Title    string
	Articles gold.Articles
	Error    error
	PrevPage int
	NextPage int
	TagCloud gold.TagCloud
	Year     int
	Month    time.Month
	Archive  []Archive
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
	Articles gold.Articles
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

func parsePage(u url.URL) int {
	if page, ok := u.Query()["page"]; ok {
		if pg, err := strconv.Atoi(page[0]); err == nil {
			return pg
		}
	}
	return 1
}

func (p *Page) MakeArchive() {
	for y, v := range data.Articles.Enabled().YearMap() {
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

func (p Page) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pg := parsePage(*r.URL)
	app := conf.Blog.ArticlesPerPage
	p.Articles, p.NextPage, p.PrevPage = p.Articles.Page(pg, app)
	p.TagCloud = data.Articles.TagCloud()
	p.Config = conf
	if p.Year == 0 {
		p.Year = time.Now().Year()
		if p.Month == 0 {
			p.Month = time.Now().Month()
		}
	}
	p.MakeArchive()

	err := tmpl.ExecuteTemplate(w, "index.tmpl", p)
	if err != nil {
		log.Fatal(err)
	}
}

type TagPage struct{ Page }

func (p *TagPage) Select(match []string) {
	s := match[0]
	p.Articles = data.Articles.Tag(s)
	p.Title = fmt.Sprint(conf.Blog.Title, " - ", s)
}

type IndexPage struct{ Page }

func (p *IndexPage) Select(match []string) {
	p.Articles = data.Articles.Enabled()
	p.Title = conf.Blog.Title
}

type SlugPage struct{ Page }

func (p *SlugPage) Select(match []string) {
	a, err := data.Articles.Find(match[0])
	if err == nil {
		p.Title = a.Title
		p.Articles = gold.Articles{a}
		p.Year = a.Year()
		p.Month = a.Month()
	} else {
		p.Title = conf.Blog.Title
	}
}

type YearPage struct{ Page }

func (p *YearPage) Select(match []string) {
	p.Year = atoiMust(match[0])
	p.Articles = data.Articles.Year(p.Year)
	p.Title = fmt.Sprint(conf.Blog.Title, " - ", p.Year)
}

type MonthPage struct{ Page }

func (p *MonthPage) Select(match []string) {
	p.Year = atoiMust(match[0])
	p.Month = time.Month(atoiMust(match[1]))
	p.Articles = data.Articles.Year(p.Year).Month(p.Month)
	p.Title = fmt.Sprint(conf.Blog.Title, " - ", p.Year, p.Month)
}