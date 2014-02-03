package helpers

import (
	"crypto/rand"
	"encoding/base64"
)

const (
	idlength = 7
)

func NewId(inuse func(string) bool) string {
	b := make([]byte, idlength)
	for {
		rand.Read(b)
		id := base64.StdEncoding.EncodeToString(b)
		if !inuse(id) {
			return id
		}
	}
}
