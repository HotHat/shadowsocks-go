package ss2022

import (
	"fmt"
	"testing"
)

func TestNonce(t *testing.T) {
	//n := EncryptNonce()
	//fmt.Println(n)

}

//var key, _ = hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")

type StructTest struct {
	private int
	Public  int
}

func TestEncryption(t *testing.T) {
	var a = StructTest{
		private: 10,
		Public:  20,
	}

	fmt.Println(a.private)
	fmt.Println(a.Public)
	a.private = 30
	a.Public = 40
	fmt.Println(a.private)
	fmt.Println(a.Public)
}

func TestDecryption(t *testing.T) {
	protocol := "2022-blake3-aes-128-gcm"
	psk := "JIloOlaO1V506UnRV521mg=="
	salt := NewSalt(CryptoGCMBlockLength(protocol))

	_ = NewCryptoGCM(NewSessionSubKey([]byte(psk), salt))

}
