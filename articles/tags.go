package articles

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

func (t TagCloud) Len() int      { return len(t) }
func (t TagCloud) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

type byWight struct{ TagCloud }

func (t byWight) Less(i, j int) bool {
	return t.TagCloud[i].Wight < t.TagCloud[j].Wight
}

type byName struct{ TagCloud }

func (t byName) Less(i, j int) bool {
	return t.TagCloud[i].Tag < t.TagCloud[j].Tag
}

func (a Articles) CountTags() TagCount {
	tags := make(TagCount)
	for _, article := range a.Enabled() {
		for _, tag := range article.Tags {
			tags[tag]++
		}
	}
	return tags
}

func (a Articles) TagMap() TagMap {
	tm := make(TagMap)
	for tag := range a.CountTags() {
		for _, article := range a.Enabled() {
			if article.Tags.Has(tag) {
				tm[tag] = append(tm[tag], article)
			}
		}
	}
	return tm
}

func (a Articles) Tag(tag string) Articles {
	var A Articles
	for _, v := range a {
		if v.Tags.Has(tag) {
			A = append(A, v)
		}
	}
	return A
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

func (a Articles) TagCloud() TagCloud {
	var tc TagCloud
	for k, v := range a.CountTags() {
		tc = append(tc, tagCloud{Tag: k, Wight: 5 / v})
	}
	sort.Sort(byName{tc})
	return tc
}
