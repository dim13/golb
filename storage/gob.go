package storage

import (
	"encoding/gob"
	"io/ioutil"
	"os"
)

func Load(fname string, v interface{}) error {
	data, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer data.Close()

	dec := gob.NewDecoder(data)
	err = dec.Decode(v)
	if err != nil {
		return err
	}

	return nil
}

func Store(fname string, v interface{}) error {
	data, err := ioutil.TempFile(os.TempDir, fname)
	if err != nil {
		return err
	}
	defer data.Close()

	enc := gob.NewEncoder(data)
	err = enc.Encode(v)
	if err != nil {
		return err
	}

	return os.Rename(data.Name(), fname)
}
