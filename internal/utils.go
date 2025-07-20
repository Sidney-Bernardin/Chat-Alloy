package internal

import (
	"crypto/rand"
	"encoding/base64"
)

func MustRandomString(len int) string {
	b := make([]byte, len)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)[:len]
}
