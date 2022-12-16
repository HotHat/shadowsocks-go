package ss2022

import (
	"crypto/cipher"
	"fmt"
)

var encryptNonce = NewCounter()

// copy of encrypt nonce syn increment
var decryptNonce = *encryptNonce

func Encryption(aead cipher.AEAD, plaintext []byte) []byte {
	nonce := encryptNonce.Bytes()
	fmt.Println("encryption nonce: ", nonce)
	ciphertext := aead.Seal(nil, nonce[:], plaintext, nil)
	encryptNonce.Increment()
	return ciphertext
}

func Decryption(aead cipher.AEAD, ciphertext []byte) ([]byte, error) {
	nonce := decryptNonce.Bytes()
	fmt.Println("decryption nonce: ", nonce)
	plaintext, err := aead.Open(nil, nonce[:], ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	decryptNonce.Increment()
	return plaintext, err
}
