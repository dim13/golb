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
	return c[i].Date.Before(c[j].Date)
}

func (c *Comments) Add(comment *Comment) {
	comment.Date = time.Now()
	*c = append(*c, comment)
}
