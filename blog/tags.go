package blog

import (
	"sort"
	"strings"
	"unicode"
)

type Tags []string
type TagMap map[string]Articles

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

func (b Blog) TagMap() TagMap {
	tm := make(TagMap)
	for _, art := range b {
		for _, tag := range art.Tags {
			tm[tag] = append(tm[tag], art)
		}
	}
	return tm
}

func (b Blog) TagCloud() (tc TagCloud) {
	for tag, art := range b.TagMap() {
		tc = append(tc, tagCloud{Tag: tag, Wight: 5 / len(art)})
	}
	sort.Sort(byName{tc})
	return
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

func (t Tags) String() string {
	return strings.Join(t, " ")
}
