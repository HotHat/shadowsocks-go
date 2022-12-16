package ss2022

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestNonce(t *testing.T) {
	//n := EncryptNonce()
	//fmt.Println(n)

}

var key, _ = hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")

func TestEncryption(t *testing.T) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	nonce := NewCounter()
	c := Crypto{aesgcm, *nonce, *nonce}

	//fmt.Println(EncryptNonce())
	b := c.Encryption([]byte("example"))

	fmt.Println(b)
	//fmt.Println(EncryptNonce())
}

func TestDecryption(t *testing.T) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := NewCounter()
	c1 := NewCrypto(aesgcm, *nonce)
	c2 := NewCrypto(aesgcm, *nonce)

	//fmt.Println("EncryptNonce:", EncryptNonce())
	b1 := c1.Encryption([]byte("example1"))
	b2 := c1.Encryption([]byte("example2"))

	p, _ := c2.Decryption(b1)
	fmt.Println(string(p))
	p, _ = c2.Decryption(b2)
	fmt.Println(string(p))

	p1 := c2.Encryption([]byte("example1111"))
	p2 := c2.Encryption([]byte("example2222"))

	q, _ := c1.Decryption(p1)
	fmt.Println(string(q))
	q, _ = c1.Decryption(p2)
	fmt.Println(string(q))

	//fmt.Println(EncryptNonce())
}
