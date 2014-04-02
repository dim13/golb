package golb

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	DataBase          string
	BlotTheme         string
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
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return &Config{}, err
	}
	err = json.Unmarshal(data, c)
	if err != nil {
		return &Config{}, err
	}
	return c, nil
}
