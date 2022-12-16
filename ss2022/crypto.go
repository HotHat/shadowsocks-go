package ss2022

import (
	"crypto/cipher"
)

type Crypto struct {
	Cipher       cipher.AEAD
	EncryptNonce Counter
	DecryptNonce Counter
}

func NewCrypto(cipher cipher.AEAD, nonce Counter) *Crypto {
	return &Crypto{
		cipher,
		nonce,
		nonce,
	}
}

func (c Crypto) Encryption(plaintext []byte) []byte {
	nonce := c.EncryptNonce.Bytes()
	ciphertext := c.Cipher.Seal(nil, nonce[:], plaintext, nil)
	c.EncryptNonce.Increment()
	return ciphertext
}

func (c Crypto) Decryption(ciphertext []byte) ([]byte, error) {
	nonce := c.DecryptNonce.Bytes()

	plaintext, err := c.Cipher.Open(nil, nonce[:], ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	c.DecryptNonce.Increment()
	return plaintext, err
}
