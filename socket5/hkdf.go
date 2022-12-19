package socket5

import (
	"crypto/sha1"
	"golang.org/x/crypto/hkdf"
	"io"
)

func NewSHA1(secret, salt, info []byte, size int) ([]byte, error) {
	out := make([]byte, size)
	hk := hkdf.New(sha1.New, secret, salt, info)
	_, err := io.ReadFull(hk, out)
	if err != nil {
		panic("hkdf sha1 fail")
	}
	return out, err
}
