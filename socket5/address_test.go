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
	b := NewSocket5DomainAddress("baidu.com", 443)
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
	b := NewSocket5DomainAddress("sina.com", 443)
	fmt.Println(b)

	addr, err := ParseSocket5Address(b)
	if err != nil {
		fmt.Println("parse error", err)
	}
	fmt.Println(addr)
	fmt.Println(addr.Length())
}
