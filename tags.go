package gold

import (
	"strings"
)

type Tags []string
type TagMap map[string]Articles
type TagCount map[string]int

func (a Articles) CountTags() TagCount {
	tags := make(TagCount)
	for _, article := range a {
		for _, tag := range article.Tags {
			tags[tag]++
		}
	}
	return tags
}

func (a *Articles) TagMap() TagMap {
	tm := make(TagMap)
	for tag, _ := range a.CountTags() {
		for _, article := range *a {
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
