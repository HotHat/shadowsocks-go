package ss2022

import (
	"crypto/rand"
	"errors"
	"github.com/zeebo/blake3"
	"log"
	util "shadowsocks-go/binary"
	"shadowsocks-go/socket5"
)

const (
	HeaderTypeClientStream = 0
	HeaderTypeServerStream = 1
	MinPaddingLength       = 0
	MaxPaddingLength       = 900
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

func NewRequestFixedHeader(typ byte, tm uint64, lg uint16) []byte {
	buf := make([]byte, 11)
	buf[0] = typ
	copy(buf[1:9], util.PutUint64(tm))
	copy(buf[9:], util.PutUint16(lg))
	return buf
}

func ParseRequestFixedHeader(buf []byte) (uint64, uint16, error) {
	if len(buf) < 11 {
		return 0, 0, errors.New("parse request fixed header error")
	}
	if buf[0] != HeaderTypeClientStream {
		return 0, 0, errors.New("client header type must be 0")
	}

	tm := util.GetUint64(buf[1:9])
	l := util.GetUint16(buf[9:])
	return tm, l, nil
}

func NewRequestVariableHeader(typ byte, addr []byte, port uint16, padLen uint16, payload []byte) []byte {
	buf := make([]byte, 0)
	var n uint16

	buf = append(buf, typ)
	n += 1

	buf = append(buf, addr...)
	n += uint16(len(addr))

	buf = append(buf, util.PutUint16(port)...)
	n += 2

	buf = append(buf, util.PutUint16(padLen)...)
	n += 2

	tmp := make([]byte, padLen)
	_, err := rand.Read(tmp)
	if err != nil {
		panic("generate rand padding content error")
	}
	buf = append(buf, tmp...)
	n += padLen

	buf = append(buf, payload...)
	n += uint16(len(payload))

	return buf
}

func ParseRequestVariableHeader(buf []byte) (*socket5.Address, []byte, error) {
	addr, err := socket5.ParseSocket5Address(buf)
	if err != nil {
		return nil, nil, err
	}
	l := addr.Length()
	if uint16(len(buf)) < l+2 {
		return nil, nil, errors.New("parse padding length error")
	}

	pl := util.GetUint16(buf[l : l+2])
	if uint16(len(buf)) < l+1+pl {
		return nil, nil, errors.New("parse padding error")
	}
	payload := buf[l+2+pl:]
	return addr, payload, nil
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
