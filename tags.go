package gold

type TagMap map[string]Articles

type Tags map[string]int

func (a Articles) CountTags() Tags {
	tags := make(Tags)
	for _, article := range a {
		for _, tag := range article.Tags {
			tags[tag]++
		}
	}
	return tags
}

func (a Article) hasTag(tag string) bool {
	for _, t := range a.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (a *Articles) TagMap() TagMap {
	tm := make(TagMap)
	for tag, _ := range a.CountTags() {
		for _, article := range *a {
			if article.hasTag(tag) {
				tm[tag] = append(tm[tag], article)
			}
		}
	}
	return tm
}
