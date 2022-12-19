package ss2022

import "crypto/cipher"

const (
	NonceLength = 12
)

type Stream struct {
	psk    []byte
	nonce  []byte
	cipher cipher.AEAD
}

func NewStream(psk []byte) Stream {
	return Stream{
		psk:   psk,
		nonce: make([]byte, NonceLength),
	}
}

func (c *Stream) AddSalt(salt []byte) {
	c.cipher = NewCryptoGCM(NewSessionSubKey(c.psk, salt))
}

func (c Stream) Encryption(plaintext []byte) []byte {
	ciphertext := c.cipher.Seal(nil, c.nonce, plaintext, nil)
	increment(c.nonce)
	return ciphertext
}

func (c Stream) Decryption(ciphertext []byte) ([]byte, error) {
	plaintext, err := c.cipher.Open(nil, c.nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	increment(c.nonce)
	return plaintext, err
}

func increment(counter []byte) {
	for i, _ := range counter {
		counter[i]++
		if counter[i] != 0 {
			break
		}
	}
}
