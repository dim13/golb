package gold

import (
	"code.google.com/p/gcfg"
)

type blog struct {
	Title    string
	Subtitle string
	Url      string
	Owner    string
	ArticlesPerPage int
	TagsInCloud     int
	DataBase        string
}

type captcha struct {
	Public  string
	Private string
}

type comments struct {
	Maxlen  int
	Allowed bool
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

func ReadConf(fname string) (*Config, error) {
	c := new(Config)
	if err := gcfg.ReadFileInto(c, fname); err != nil {
		return &Config{}, err
	}
	return c, nil
}
