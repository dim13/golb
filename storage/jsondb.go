package storage

import (
	"encoding/json"
	"io/ioutil"
)

func ReadJson(fname string, v interface{}) error {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func WriteJson(fname string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fname, data, 0664)
}
