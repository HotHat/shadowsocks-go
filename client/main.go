package main

import (
	"context"
	"fmt"
	"net"
	"shadowsocks-go/websocket"
	"time"
)

func main() {
	listen, err := net.Listen("tcp", "0.0.0.0:9000")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept fail:", err)
			continue
		}

		fmt.Println("Connect from", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handle(ctx context.Context, duration time.Duration) {
	ctx1, cancel := context.WithCancel(ctx)

	fmt.Println("after block")
	cancel()

	select {
	case <-ctx.Done():
		fmt.Println("handle", ctx1.Err())
	case <-time.After(duration):
		fmt.Println("process request with", duration)
	}

}

func handleRead(ctx context.Context, conn net.Conn, channel chan<- []byte, isExit chan<- bool) {
	buf := make([]byte, 4096)
	readLen := 0
	for {
		tmp := buf[readLen:]
		if len(tmp) == 0 {
			isExit <- true
			break
		}

		n, err := conn.Read(tmp)
		readLen += n
		if err != nil {
			isExit <- true
			break
		}
		fmt.Println("one time read:", string(tmp[:n]))
		fmt.Println("one request read:", string(buf[:readLen]))

		// parse http header
		h, ln, err1 := websocket.ParseHttpHeaders(buf[:readLen])
		if err1 != nil {
			if err1.IsContinue() {
				continue
			}
			fmt.Println(err1)
			isExit <- true
			break
		}

		fmt.Println("http header:", h)
		fmt.Println("http header len:", ln)
		channel <- buf[:readLen]

		// some data left in buffer
		if ln < readLen {
			copy(buf, buf[ln:readLen])
			readLen = readLen - ln
			fmt.Println("data left in buffer:", buf[:readLen])
		} else {
			// 重新开始
			readLen = 0
		}

		//
		break
	}

	for {
		tmp := buf[readLen:]
		if len(tmp) == 0 {
			isExit <- true
			break
		}

		n, err := conn.Read(tmp)
		readLen += n
		if err != nil {
			isExit <- true
			break
		}

		fin, op, mask, pl, hl, err1 := websocket.ParseFramePayloadLength(buf[:readLen])

		if err1 != nil {
			if err1.IsContinue() {
				continue
			}
		}

		dataBuf := make([]byte, pl)
		// data in buf
		if int(hl) < readLen {
			copy(dataBuf, buf[hl:readLen])
		}

		readLen = 0
		for readLen < int(pl) {
			tp := dataBuf[readLen:]
			n, err := conn.Read(tp)
			readLen += n
			if err != nil {
				isExit <- true
				break
			}
		}

		// dataBuf with the frame data
		if len(mask) > 0 {
			for k, _ := range dataBuf {
				m := mask[k%4]
				dataBuf[k] ^= m
			}
		}

		fmt.Printf("fin: %b\n", fin)
		fmt.Printf("op:  %b\n", op)
		fmt.Printf("data: %s\n", string(dataBuf))
	}

	select {
	case <-ctx.Done():
		return
	}
}

func handleWrite(ctx context.Context, conn net.Conn, channel <-chan []byte, isExit chan<- bool) {
	for {
		select {
		case buf := <-channel:
			n, err := conn.Write(buf)
			if err != nil {
				isExit <- true
			}
			fmt.Println("Write:", n, " data:", buf)
		case <-ctx.Done():
			return
		}
	}
}

func handleConnection(conn net.Conn) {
	ctx, cancel := context.WithCancel(context.Background())
	isExit := make(chan bool)
	channel := make(chan []byte)
	defer close(isExit)
	defer close(channel)

	go handleRead(ctx, conn, channel, isExit)
	go handleWrite(ctx, conn, channel, isExit)

	for {
		select {
		case <-isExit:
			{
				cancel()
				conn.Close()
				fmt.Println("close connect from", conn.RemoteAddr())
			}
		}
	}

}
