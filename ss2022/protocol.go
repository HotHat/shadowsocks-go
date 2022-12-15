package ss2022

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/zeebo/blake3"
	"log"
	"net"
	"time"
)

const (
	HeaderTypeClientStream = 0
	HeaderTypeServerStream = 1
	MinPaddingLength       = 0
	MaxPaddingLength       = 900

	Socket5AddressIPv4   = 0x01
	Socket5AddressDomain = 0x03
	Socket5AddressIPv6   = 0x04
)

type RequestFixedHeader struct {
	Type      byte
	Timestamp uint64
	Length    uint16
}

type RequestVariableHeader struct {
	AType         byte
	Address       []byte
	Port          uint16
	PaddingLength uint16
	Payload       []byte
}

type ResponseFixedHeader struct {
	Type        byte
	Timestamp   uint64
	RequestSalt []byte
	Length      uint16
}

func (r *RequestFixedHeader) New(l uint16) {
	r.Type = 1
	r.Timestamp = Timestamp()
	r.Length = l
}

func (r *RequestFixedHeader) Bytes() []byte {
	buf := make([]byte, 11)
	buf[0] = r.Type
	copy(buf[1:9], PutUint64(r.Timestamp))
	copy(buf[9:], PutUint16(r.Length))

	return buf
}

func (r *RequestVariableHeader) New(atyp byte, addr []byte, port uint16, padLen uint16, payload []byte) {
	r.AType = atyp
	r.Address = addr
	r.Port = port
	r.PaddingLength = padLen
	r.Payload = payload
}

func (r *RequestVariableHeader) Bytes() []byte {
	buf := make([]byte, 0)
	var n uint16

	buf = append(buf, r.AType)
	n += 1

	buf = append(buf, r.Address...)
	n += uint16(len(r.Address))

	buf = append(buf, PutUint16(r.Port)...)
	n += 2

	buf = append(buf, PutUint16(r.PaddingLength)...)
	n += 2

	tmp := make([]byte, r.PaddingLength)
	buf = append(buf, tmp...)
	n += r.PaddingLength

	buf = append(buf, r.Payload...)
	n += uint16(len(r.Payload))

	return buf
}

func NewSocket5IPv4Address(ip net.IP, port uint16) []byte {
	buf := make([]byte, 1+4+2)

	buf[0] = Socket5AddressIPv4

	copy(buf[1:], ip.To4())
	copy(buf[5:], PutUint16(port))

	return buf
}

func NewSocket5IPv6Address(ip net.IP, port uint16) []byte {
	buf := make([]byte, 1+16+2)

	buf[0] = Socket5AddressIPv4

	copy(buf[1:], ip.To16())
	copy(buf[17:], PutUint16(port))

	return buf
}

func NewSocket5DomainAddress(domain string, port uint16) []byte {
	buf := make([]byte, 1+2+256)

	buf[0] = Socket5AddressDomain

	c := []byte(domain)
	l := len(c)
	copy(buf[1:], c)
	copy(buf[l+1:], PutUint16(port))

	return buf[0 : l+3]
}

type RequestStream struct {
	Salt      []byte
	FixHeader []byte
	Chunk     []byte
}

type ResponseStream struct {
	Salt      []byte
	FixHeader []byte
	Chunk     []byte
}

func NewSalt(len int) []byte {
	if len != 16 && len != 32 {
		panic("salt length must be 16 or 32")
	}
	b := make([]byte, len)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalln(err)
	}
	return b
}

func NewSessionSubKey(context string, key []byte, salt []byte) []byte {
	material := key[:]
	copy(material, salt)
	out := make([]byte, 32)
	blake3.DeriveKey(context, material, out)
	return out
}

func Base64Decoder(encode string) []byte {
	buf := make([]byte, 512)
	r := bytes.NewReader([]byte(encode))
	decoder := base64.NewDecoder(base64.StdEncoding, r)
	n, err := decoder.Read(buf)
	if err != nil {
		log.Println(err)
	}

	return buf[0:n]
}

func Timestamp() uint64 {
	now := time.Now()
	fmt.Println(now)

	tm := now.Unix()
	return uint64(tm)
}

// PutUint16 convert uint16 to big endian byte
func PutUint16(num uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, num)
	return b
}

// PutUint64 convert uint64 to big endian byte
func PutUint64(num uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, num)
	return b
}
