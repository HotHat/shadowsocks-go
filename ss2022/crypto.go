package ss2022

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/zeebo/blake3"
)

func NewSessionSubKey(key []byte, salt []byte) []byte {
	material := key[:]
	copy(material, salt)
	out := make([]byte, 32)
	blake3.DeriveKey("shadowsocks 2022 session subkey", material, out)
	return out
}

func NewCryptoGCM(key []byte) cipher.AEAD {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	return aesgcm
}

func CryptoGCMBlockLength(pro string) int {
	if pro == "2022-blake3-aes-128-gcm" {
		return 16
	} else if pro == "2022-blake3-aes-256-gcm" {
		return 32
	}

	return 0
}
