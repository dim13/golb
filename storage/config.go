package storage

import (
	"log"

	"code.google.com/p/gcfg"
)

type Blog struct {
	Title           string
	Description     string
	Owner           string
	ArticlesPerPage int
	DataBase        string
}

type Captcha struct {
	Public  string
	Private string
}

type Comments struct {
	Maxlen  int
	Enabled bool
}

type SMTP struct {
	Server string
	Sender string
}

type Google struct {
	AnalyticsID string
	WebmasterID string
}

type Config struct {
	Blog     Blog
	Captcha  Captcha
	Comments Comments
	Smtp     SMTP
	Google   Google
}

func ReadConf(fname string) (c Config, err error) {
	log.Println("Read", fname)
	return c, gcfg.ReadFileInto(&c, fname)
}
