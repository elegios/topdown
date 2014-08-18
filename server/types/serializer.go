package types

import (
	"encoding/gob"
	"os"
	"path/filepath"
)

type Decoder func(path string, v interface{}) error
type Encoder func(path string, v interface{}) error

var (
	enc Encoder = GobEncode
	dec Decoder = GobDecode
)

func init() {
	gob.Register(Position{})
}

func GobEncode(path string, v interface{}) (err error) {
	err = os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return
	}
	fi, err := os.Create(path)
	if err != nil {
		return
	}
	defer fi.Close()

	return gob.NewEncoder(fi).Encode(v)
}

func GobDecode(path string, v interface{}) (err error) {
	fi, err := os.Open(path)
	if err != nil {
		return
	}
	defer fi.Close()

	return gob.NewDecoder(fi).Decode(v)
}
