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

func EvpBytesToKey(password string, keyLen int) (key []byte) {
	const md5Len = 16

	cnt := (keyLen-1)/md5Len + 1
	m := make([]byte, cnt*md5Len)
	copy(m, MD5([]byte(password)))
	d := make([]byte, md5Len+len(password))
	start := 0
	for i := 1; i < cnt; i++ {
		start += md5Len
		copy(d, m[start-md5Len:start])
		copy(d[md5Len:], password)
		copy(m[start:], MD5(d))
	}
	return m[:keyLen]
}

func MD5(data []byte) []byte {
	hash := md5.New()
	hash.Write(data)
	return hash.Sum(nil)
}
