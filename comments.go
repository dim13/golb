// Comments
package golb

import (
	"time"
)

type Comments []*Comment

type Comment struct {
	Date    time.Time
	Name    string
	Email   string
	URL     string
	Comment string
	Enabled bool
}

func (c *Comment) Publish() {
	c.Enabled = true
}

func (c *Comment) Suppress() {
	c.Enabled = false
}

func (c Comments) Len() int {
	return len(c)
}

func (c Comments) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c Comments) Less(i, j int) bool {
	return c[i].Date.After(c[j].Date)
}

func (cp Comments) Add(c *Comment) Comments {
	c.Date = time.Now()
	return append(cp, c)
}
