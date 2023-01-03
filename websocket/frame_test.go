package websocket

import (
	"encoding/base64"
	"fmt"
	"net"
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

func TestNewFrame2(t *testing.T) {

	up := NewHttpUpgrade("/", "baidu.com", "baidu.com")
	s := NewFrame(true, OpcodePing, false, []byte("hello, websocket"))
	//s := []byte{0x81, 0x85, 0x37, 0xfa, 0x21, 0x3d, 0x7f, 0x9f, 0x4d, 0x51, 0x58}

	//b := s[6:]

	//fmt.Printf("%b", s)
	//fmt.Println(len(b))
	//fmt.Println(b)
	//fmt.Println(string(b))
	//_, _, p, _ := ParseFrame(s)
	//fmt.Println(string(p))
	///*
	conn, err := net.Dial("tcp", "192.168.33.10:8088")
	if err != nil {
		panic(err)
	}
	buff := make([]byte, 4096)

	n, err := conn.Write(up)
	if err != nil {
		panic(err)
	}
	fmt.Println("write byte:", n)

	n, _ = conn.Read(buff)
	fmt.Println(string(buff[:n]))

	fmt.Printf("%b", s)
	n, err = conn.Write(s)
	if err != nil {
		panic(err)
	}
	fmt.Println("write byte:", n)

	n, _ = conn.Read(buff)
	fin, opcode, p, _ := ParseFrame(buff[:n])
	fmt.Println("fin:", fin)
	fmt.Println("opcode:", opcode)
	fmt.Println("receive:", string(p))

	//*/

}

func TestEn(t *testing.T) {
	s := base64.StdEncoding.EncodeToString([]byte{202, 100, 234, 10})
	fmt.Println(s)
}

func TestParseFrame2(t *testing.T) {
	b := []byte{0x81, 0x05, 0x48, 0x65, 0x6c, 0x6c, 0x6f}
	f, op, p, _ := ParseFrame(b)
	fmt.Println(f)
	fmt.Println(op)
	fmt.Println(string(p))
}

func TestParseFrame3(t *testing.T) {
	b := []byte{0x81, 0x85, 0x37, 0xfa, 0x21, 0x3d, 0x7f, 0x9f, 0x4d, 0x51, 0x58}
	f, op, p, _ := ParseFrame(b)
	fmt.Println(f)
	fmt.Println(op)
	fmt.Println(string(p))
}

func TestParseFramePayloadLengthWithoutMask(t *testing.T) {
	b := []byte{0x81, 0x05, 0x48, 0x65, 0x6c, 0x6c, 0x6f}
	f, op, m, p, hl, err := ParseFramePayloadLength(b)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(f)
	fmt.Println(op)
	fmt.Println(m)
	fmt.Println(p)
	fmt.Println(hl)
}

func TestParseFramePayloadLengthWithMask(t *testing.T) {
	b := []byte{0x81, 0x85, 0x37, 0xfa, 0x21, 0x3d, 0x7f, 0x9f, 0x4d, 0x51, 0x58}
	f, op, m, p, hl, err := ParseFramePayloadLength(b)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(f)
	fmt.Println(op)
	fmt.Println(m)
	fmt.Println(p)
	fmt.Println(hl)
}

func TestHttpUpgradeKeyValidate(t *testing.T) {
	b := HttpUpgradeKeyValidate("dGhlIHNhbXBsZSBub25jZQ==", "s3pPLMBiTxaQ9kYGzzhZRbK+xOo=")
	fmt.Println(b)
}

func TestParseHttpHeaders(t *testing.T) {
	str := "GET /chat HTTP/1.1\r\nHost: server.example.com\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==\r\nOrigin: http://example.com\r\nSec-WebSocket-Protocol: chat, superchat\r\nSec-WebSocket-Version: 13\r\n"
	//str := "GET /chat HTTP/1.1\r\n"

	// not found http header end
	h, e, err := ParseHttpHeaders([]byte(str))

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("map:", h)
		fmt.Println("end:", e)
	}

	//  normal http header
	str += "\r\n"
	h, e, err = ParseHttpHeaders([]byte(str))

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("map:", h)
		fmt.Println("end:", e, " string len:", len(str))
	}

}
