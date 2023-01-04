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

	b := NewFrame(true, OpcodeText, true, []byte("Read reads data from the connection. Read can be made to time out and return an error after a fixed time limit; see SetDeadline and SetReadDeadline"))
	conn.Write(b)

	b = NewFrame(true, OpcodeText, true, []byte("ReadSetReadDeadline"))
	conn.Write(b)

	b = NewFrame(true, OpcodeText, false, []byte("LocalAddr returns the local network address, if known."))
	conn.Write(b)

	//for {
	time.Sleep(5000)
	//}

	conn.Close()

}
