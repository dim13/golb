package gold

import (
	"sort"
	"strings"
)

type Tags []string
type TagMap map[string]Articles
type TagCount map[string]int

type TagCloud []tagCloud
type tagCloud struct {
	Tag   string
	Wight int
}

type ByWight []tagCloud

func (t ByWight) Len() int           { return len(t) }
func (t ByWight) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByWight) Less(i, j int) bool { return t[i].Wight < t[j].Wight }

type ByName []tagCloud

func (t ByName) Len() int           { return len(t) }
func (t ByName) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByName) Less(i, j int) bool { return t[i].Tag < t[j].Tag }

func (a Articles) CountTags() TagCount {
	tags := make(TagCount)
	for _, article := range a {
		for _, tag := range article.Tags {
			tags[tag]++
		}
	}
	return tags
}

func (a Articles) TagMap() TagMap {
	tm := make(TagMap)
	for tag, _ := range a.CountTags() {
		for _, article := range a {
			if article.Tags.Has(tag) {
				tm[tag] = append(tm[tag], article)
			}
		}
	}
	return tm
}

func (ts Tags) Has(tag string) bool {
	for _, t := range ts {
		if t == tag {
			return true
		}
	}
	return false
}

func (t Tags) String() string {
	return strings.Join(t, ",")
}

func ReadTags(s string) Tags {
	return strings.Split(s, ",")
}

func (a Articles) TagCloud(n int) TagCloud {
	var tc TagCloud
	for k, v := range a.CountTags() {
		tc = append(tc, tagCloud{Tag: k, Wight: v})
	}
	sort.Sort(sort.Reverse(ByWight(tc)))
	tc = tc[:n]
	bottom := tc[len(tc)-1].Wight
	for i, _ := range tc {
		tc[i].Wight -= bottom - 1
		tc[i].Wight = 5 - 5 / tc[i].Wight
	}
	sort.Sort(ByName(tc))
	return tc
}
