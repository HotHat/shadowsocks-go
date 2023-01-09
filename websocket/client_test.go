package websocket

import (
	"fmt"
	"net"
	"net/url"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	conn, err := net.Dial("tcp", "0.0.0.0:9000")
	if err != nil {
		panic(err)
	}

	httpHeader := "GET /chat HTTP/1.1\r\nHost: server.example.com\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==\r\nOrigin: http://example.com\r\nSec-WebSocket-Protocol: chat, superchat\r\nSec-WebSocket-Version: 13\r\n\r\n"
	conn.Write([]byte(httpHeader))

	/*
		b := NewFrame(true, OpcodeText, true, []byte("Read reads data from the connection. Read can be made to time out and return an error after a fixed time limit; see SetDeadline and SetReadDeadline"))
		conn.Write(b)

		b = NewFrame(true, OpcodeText, true, []byte("ReadSetReadDeadline"))
		conn.Write(b)

		b = NewFrame(true, OpcodeText, false, []byte("LocalAddr returns the local network address, if known."))
		conn.Write(b)
	*/
	conn.Write([]byte{0x01, 0x03, 0x48, 0x65, 0x6c})
	conn.Write([]byte{0x80, 0x02, 0x6c, 0x6f})

	//for {
	time.Sleep(5000)
	//}

	conn.Close()

}

func TestUrl(t *testing.T) {
	p, err := url.Parse("ws://game.example.com:12010/updates")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p)
	fmt.Println(p.Host)
	fmt.Println(p.Scheme)
	fmt.Println(p.Path)
	fmt.Println(p.Port())
}
