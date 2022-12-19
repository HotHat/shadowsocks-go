package utils

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"golang.org/x/crypto/hkdf"
	"io"
	"log"
)

func ToByte(s string) []byte {
	return []byte(s)
}

func NewSalt(len int) []byte {
	if len != 16 && len != 32 {
		panic("salt length must be 16 or 32")
	}
	b := make([]byte, len)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalln(err)
	}
	return b
}

func IncrementNonce(nonce []byte) {
	for i, _ := range nonce {
		nonce[i]++
		if nonce[i] != 0 {
			break
		}
	}
}

func Kdf(password string, keyLen int) []byte {
	var b, prev []byte
	h := md5.New()
	for len(b) < keyLen {
		h.Write(prev)
		h.Write([]byte(password))
		b = h.Sum(b)
		prev = b[len(b)-h.Size():]
		h.Reset()
	}
	return b[:keyLen]
}
func KdfSHA1(secret, salt, info, out []byte) {
	hk := hkdf.New(sha1.New, secret, salt, info)
	_, err := io.ReadFull(hk, out)
	if err != nil {
		panic("hkdf sha1 fail")
	}
}
