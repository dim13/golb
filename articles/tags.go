package articles

import (
	"sort"
	"strings"
	"unicode"
)

type Tags []string
type TagMap map[string]Articles
type tagCount map[string]int

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

func (a Articles) countTags() tagCount {
	tags := make(tagCount)
	for _, article := range a {
		for _, tag := range article.Tags {
			tags[tag]++
		}
	}
	return tags
}

func (a Articles) TagMap() TagMap {
	tm := make(TagMap)
	for tag := range a.countTags() {
		for _, article := range a {
			if article.Tags.Has(tag) {
				tm[tag] = append(tm[tag], article)
			}
		}
	}
	return tm
}

func (a Articles) Tag(tag string) (A Articles) {
	for _, v := range a {
		if v.Tags.Has(tag) {
			A = append(A, v)
		}
	}
	return
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
	return strings.Join(t, " ")
}

func uniq(in []string) (out []string) {
	m := make(map[string]struct{})
	for _, v := range in {
		m[v] = struct{}{}
	}
	for k := range m {
		out = append(out, k)
	}
	sort.Sort(sort.StringSlice(out))
	return
}

func ReadTags(s string) Tags {
	f := func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsPunct(r)
	}
	return uniq(strings.FieldsFunc(s, f))
}

func (a Articles) TagCloud() (tc TagCloud) {
	for k, v := range a.countTags() {
		tc = append(tc, tagCloud{Tag: k, Wight: 5 / v})
	}
	sort.Sort(byName{tc})
	return
}
