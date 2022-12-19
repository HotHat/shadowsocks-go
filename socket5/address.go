package socket5

import (
	"errors"
	"fmt"
	"net"
)
import "shadowsocks-go/binary"

const (
	Socket5AddressIPv4   = 0x01
	Socket5AddressDomain = 0x03
	Socket5AddressIPv6   = 0x04
)

type Address struct {
	Type   byte
	Ip     net.IP
	Domain string
	Port   uint16
}

func (a Address) Length() uint16 {
	if a.Type == 1 {
		return 7
	}
	if a.Type == 4 {
		return 19
	}
	if a.Type == 3 {
		return uint16(2 + len(a.Domain) + 2)
	}
	return 0
}

func NewSocket5IPv4Address(ip net.IP, port uint16) []byte {
	buf := make([]byte, 1+4+2)

	buf[0] = Socket5AddressIPv4

	copy(buf[1:], ip.To4())
	copy(buf[5:], binary.PutUint16(port))

	return buf
}

func NewSocket5IPv6Address(ip net.IP, port uint16) []byte {
	buf := make([]byte, 1+16+2)

	buf[0] = Socket5AddressIPv6

	copy(buf[1:], ip.To16())
	copy(buf[17:], binary.PutUint16(port))

	return buf
}

func NewSocket5DomainAddress(domain []byte, port uint16) []byte {
	buf := make([]byte, 1+1+256+2)

	buf[0] = Socket5AddressDomain

	//c := domain
	l := len(domain)
	if l > 255 {
		panic("domain length more than 255")
	}
	buf[1] = uint8(l)
	copy(buf[2:], domain)
	copy(buf[l+2:], binary.PutUint16(port))

	return buf[0 : l+4]
}

func ParseSocket5Address(buf []byte) (*Address, error) {
	if len(buf) < 2 {
		return nil, errors.New("can't parse socket5 address")
	}
	var t byte
	var ip net.IP
	var port uint16
	var domain = ""

	t = buf[0]
	if t == 1 {
		if len(buf) < 7 {
			return nil, errors.New("can't parse socket5 IPv4 address")
		}
		ip = net.IPv4(buf[1], buf[2], buf[3], buf[4])
		port = binary.GetUint16(buf[5:])
	} else if t == 4 {
		if len(buf) < 1+16+2 {
			return nil, errors.New("can't parse socket5 IPv6 address")
		}
		ip = net.IP(buf[1:17])
		port = binary.GetUint16(buf[17:])
	} else if t == 3 {
		l := int(buf[1])
		if len(buf) < (1 + 1 + l + 2) {
			return nil, errors.New("can't parse socket5 domain address")
		}
		domain = string(buf[2 : 2+l])
		port = binary.GetUint16(buf[2+l:])
	} else {
		panic(fmt.Sprintf("socket5 not support address type %d", t))
	}

	return &Address{
		Type:   t,
		Ip:     ip,
		Domain: domain,
		Port:   port,
	}, nil
}
