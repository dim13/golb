package storage

import (
	"encoding/gob"
	"os"
)

func load(fname string, v interface{}) error {
	r, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer r.Close()
	dec := gob.NewDecoder(r)
	err = dec.Decode(v)
	if err != nil {
		return err
	}
	return nil
}

func store(fname string, v interface{}) error {
	tmpfile := fname + ".tmp"
	w, err := os.Create(tmpfile)
	if err != nil {
		return err
	}
	defer w.Close()
	enc := gob.NewEncoder(w)
	err = enc.Encode(v)
	if err != nil {
		os.Remove(tmpfile)
		return err
	}
	os.Rename(tmpfile, fname)
	return nil
}
