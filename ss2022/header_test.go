package ss2022

import (
	"fmt"
	"net"
	"shadowsocks-go/binary"
	"testing"
)

func TestNewRequestFixedHeader(t *testing.T) {
	b := NewRequestFixedHeader(0, binary.Timestamp(), 256)
	fmt.Println(b)
}

func TestParseRequestFixedHeader(t *testing.T) {
	b := NewRequestFixedHeader(0, binary.Timestamp(), 206)
	fmt.Println(b)

	tm, length, err := ParseRequestFixedHeader(b)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("timestamp:", tm)
	fmt.Println("length", length)
}

func TestNewRequestVariableHeader(t *testing.T) {
	b := NewRequestVariableHeader(1, net.IPv4(192, 168, 0, 128), 80, 10, []byte("abcdefg"))
	fmt.Println(b)
}

func TestParseRequestVariableHeader(t *testing.T) {
	b := NewRequestVariableHeader(1, net.IPv4(192, 168, 0, 128), 80, 10, []byte("abcdefg"))
	fmt.Println(b)

	addr, payload, err := ParseRequestVariableHeader(b)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("address:", addr)
	fmt.Println("payload", payload)
}
