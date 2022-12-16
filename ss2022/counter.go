package ss2022

import "crypto/rand"

type Counter struct {
	Number [12]byte
}

func NewCounter() *Counter {
	counter := new(Counter)
	_, err := rand.Read(counter.Number[:])
	if err != nil {
		panic("nonce encryptNonce init fail")
	}

	return counter
}

func (c *Counter) Increment() {

	for i, _ := range c.Number {
		c.Number[i]++
		if c.Number[i] != 0 || i == 11 {
			break
		}
	}
}

func (c *Counter) Bytes() [12]byte {
	return c.Number
}
