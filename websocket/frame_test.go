package websocket

import (
	"fmt"
	"testing"
)

func TestNewFrame1(t *testing.T) {
	data := "This wire format for the data transfer part is described by the ABNF"
	b := NewFrame(true, OpcodePing, true, []byte(data))

	fmt.Printf("%v\n", []byte(data))
	fmt.Printf("%v\n", b)
	fmt.Printf("%b\n", b[:10])
}

func TestMask(t *testing.T) {

	a1 := 84 ^ 228
	a2 := 104 ^ 166
	a3 := 105 ^ 36
	a4 := 115 ^ 40
	fmt.Printf("%b, %d\n", a1, a1)
	fmt.Printf("%b, %d\n", a2, a2)
	fmt.Printf("%b, %d\n", a3, a3)
	fmt.Printf("%b, %d\n", a4, a4)

	fmt.Printf("%d, %d\n", 84, a1^228)
	fmt.Printf("%d, %d\n", 104, a2^166)
	fmt.Printf("%d, %d\n", 105, a3^36)
	fmt.Printf("%d, %d\n", 115, a4^40)
}

func TestParseFrame(t *testing.T) {
	data := "This wire format for the data transfer part is described by the ABNF"
	b := NewFrame(true, OpcodePing, true, []byte(data))

	fin, opcode, payload, err := ParseFrame(b)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(fin)
	fmt.Println(opcode)
	fmt.Println(string(payload))
}
