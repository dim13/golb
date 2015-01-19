package main

import (
	"sort"
	"time"

	"github.com/dim13/gold/blog"
)

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

func (p *page) MakeArchive() {
	for y, v := range Blog.Articles().YearMap() {
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
