package storage

import (
	"log"

	"code.google.com/p/gcfg"
)

type blog struct {
	Title           string
	Description     string
	Owner           string
	ArticlesPerPage int
	DataBase        string
}

type captcha struct {
	Public  string
	Private string
}

type comments struct {
	Maxlen  int
	Enabled bool
}

type smtp struct {
	Server string
	Sender string
}

type google struct {
	AnalyticsID string
	WebmasterID string
}

type Config struct {
	Blog     blog
	Captcha  captcha
	Comments comments
	Smtp     smtp
	Google   google
}

func ReadConf(fname string) (c Config, err error) {
	log.Println("Read", fname)
	return c, gcfg.ReadFileInto(&c, fname)
}
