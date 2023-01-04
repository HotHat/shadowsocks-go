package websocket

import (
	"net"
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

	b := []byte{0x81, 0x85, 0x37, 0xfa, 0x21, 0x3d, 0x7f, 0x9f, 0x4d, 0x51, 0x58}
	conn.Write(b)

	for {
		time.Sleep(2000)
	}

	conn.Close()

}
