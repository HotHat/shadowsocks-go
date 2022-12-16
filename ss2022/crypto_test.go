package ss2022

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestNonce(t *testing.T) {
	//n := Nonce()
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

	//fmt.Println(Nonce())
	b := Encryption(aesgcm, []byte("example"))

	fmt.Println(b)
	//fmt.Println(Nonce())
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
	//fmt.Println("Nonce:", Nonce())
	b := Encryption(aesgcm, []byte("example"))
	fmt.Println("Encryption:", b)
	//fmt.Println("Nonce:", Nonce())

	p, err := Decryption(aesgcm, b)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println(string(p))
	//fmt.Println(Nonce())
}
