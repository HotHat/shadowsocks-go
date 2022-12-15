package ss2022

type Counter struct {
	Number [12]byte
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
