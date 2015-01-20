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

func makeArchive(yr int, mth time.Month) (A []year) {
	ym := Blog.Articles().YearMap()
	for k, v := range ym {
		year := year{
			Year:  k,
			Count: len(v),
		}
		if k == yr {
			mm := v.MonthMap()
			for k, v := range mm {
				month := month{
					Year:  yr,
					Month: time.Month(k),
					Count: len(v),
				}
				if time.Month(k) == mth {
					month.Articles = v
				}
				year.Month = append(year.Month, month)
			}
			sort.Sort(sort.Reverse(byMonth(year.Month)))
		}
		A = append(A, year)
	}
	sort.Sort(sort.Reverse(byYear(A)))
	return A
}
