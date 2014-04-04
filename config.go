package gold

import (
	"code.google.com/p/gcfg"
)

type Config struct {
	Main
}

type Main struct {
	DataBase          string
	BlogTheme         string
	BlogTitle         string
	BlogSubtitle      string
	BlogUrl           string
	BlogOwner         string
	BlogRights        string
	CaptchaPublic     string
	CaptchaPrivate    string
	CommentLength     int
	CommentsAllowed   bool
	SmtpServer        string
	SmtpSender        string
	ArticlesPerPage   int
	GoogleAnalyticsId string
	GoogleWebmasterId string
	TagsInCloud       int
}

func ReadConf(fname string) (*Config, error) {
	c := new(Config)
	if err := gcfg.ReadFileInto(c, fname); err != nil {
		return &Config{}, err
	}
	return c, nil
}
