package socket5

import (
	"fmt"
	"net"
	"testing"
)

func TestNewSocket5IPv4Address(t *testing.T) {
	b := NewSocket5IPv4Address(net.IPv4(192, 168, 1, 100), 8080)
	fmt.Println(b)
}

func TestNewSocket5IPv6Address(t *testing.T) {
	b := NewSocket5IPv6Address(net.ParseIP("2001:db8::68"), 8080)
	fmt.Println(b)
}

func TestNewSocket5DomainAddress(t *testing.T) {
	b := NewSocket5DomainAddress([]byte("baidu.com"), 443)
	fmt.Println(b)
}

func TestParseSocket5Address1(t *testing.T) {
	b := NewSocket5IPv4Address(net.IPv4(192, 168, 1, 100), 8080)
	fmt.Println(b)

	addr, err := ParseSocket5Address(b)
	if err != nil {
		fmt.Println("parse error", err)
	}
	fmt.Println(addr)
}

func TestParseSocket5Address2(t *testing.T) {
	b := NewSocket5IPv6Address(net.ParseIP("2001:db8::68"), 8080)
	fmt.Println(b)

	addr, err := ParseSocket5Address(b)
	if err != nil {
		fmt.Println("parse error", err)
	}
	fmt.Println(addr)
}

func TestParseSocket5Address3(t *testing.T) {
	b := NewSocket5DomainAddress([]byte("sina.com"), 443)
	fmt.Println(b)

	addr, err := ParseSocket5Address(b)
	if err != nil {
		fmt.Println("parse error", err)
	}
	fmt.Println(addr)
	fmt.Println(addr.Length())
}

func TestParseDomain(t *testing.T) {
	b := []byte{3, 98, 97, 105, 100, 117, 46, 99, 111, 109, 0, 80, 0, 10, 180, 8, 185, 58, 158, 174, 166, 152, 76, 93}

	addr, err := ParseSocket5Address(b)
	if err != nil {
		fmt.Println("parse error", err)
	}
	fmt.Println(addr)
	fmt.Println(addr.Length())

}
