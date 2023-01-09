package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"shadowsocks-go/websocket"
	"strings"
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
			fmt.Println("Parse Http Headers error")
			fmt.Println(err1)
			isExit <- true
			break
		}

		fmt.Println("http header:", h)
		fmt.Println("http header len:", ln)

		//channel <- buf[:readLen]

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

	isFrameStart := readLen == 0
	packBuf := bytes.Buffer{}

	for {
		fmt.Println("-------------ws frame loop------------")
		// get more buffer data
		if isFrameStart {
			tmp := buf[readLen:]
			if len(tmp) == 0 {
				goto IOError
			}

			n, err := conn.Read(tmp)
			readLen += n
			if err != nil {
				goto IOError
			}
		}

		isFrameStart = true

		fin, op, mask, pl, hl, err1 := websocket.ParseFramePayloadLength(buf[:readLen])
		if err1 != nil {
			if err1.IsContinue() {
				fmt.Println("need more frame data")
				continue
			}
		}

		fmt.Printf("fin: %b\n", fin)
		fmt.Printf("op:  %b\n", op)
		fmt.Printf("mask:  %v\n", mask)
		fmt.Printf("frame data length:  %d\n", pl)
		fmt.Printf("frame head length:  %d\n", hl)
		fmt.Printf("readLen:  %d\n", readLen)

		dataBuf := make([]byte, pl)
		idx := int(uint64(hl) + pl)
		if idx <= readLen {
			copy(dataBuf, buf[hl:idx])
			copy(buf, buf[idx:readLen])
			readLen = readLen - idx
		} else {
			copy(dataBuf, buf[hl:readLen])
			r := readLen - int(hl)
			readLen = 0

			// read more data from conn
			for r < int(pl) {
				tp := dataBuf[r:]
				n, err := conn.Read(tp)
				r += n
				if err != nil {
					goto IOError
				}
			}
		}

		// dataBuf with the frame data
		if len(mask) > 0 {
			fmt.Printf("mask: %d %v\n", len(dataBuf), dataBuf)
			for k, _ := range dataBuf {
				m := mask[k%4]
				dataBuf[k] ^= m
			}
			fmt.Printf("unmask: %v\n", string(dataBuf))
		}

		isFrameStart = readLen == 0

		packBuf.Write(dataBuf)
		if fin {
			fmt.Println("This is fin package")
			fmt.Println(string(packBuf.Bytes()))
			packBuf.Reset()
		} else {
			fmt.Println("This is not fin package")
			fmt.Println(string(packBuf.Bytes()))
		}

		fmt.Printf("data: %s\n", string(dataBuf))
	}

IOError:
	isExit <- true

	select {
	case <-ctx.Done():
		return
	}
}

func handleWrite(ctx context.Context, conn net.Conn, channel <-chan []byte, isExit chan<- bool) {
	for {
		select {
		case buf := <-channel:
			b := websocket.NewFrame(true, websocket.OpcodeText, true, buf)
			n, err := conn.Write(b)
			if err != nil {
				isExit <- true
			}
			fmt.Println("Write:", n, " data:", buf)
		case <-ctx.Done():
			return
		}
	}
}
func readFromTerminate(ctx context.Context, conn net.Conn, channel chan<- []byte, isExit chan<- bool) {
	fmt.Println("Input from terminate")

	for {
		rd := bufio.NewReader(os.Stdin)
		lineBuf, _, err := rd.ReadLine()
		if err != nil {
			fmt.Println(err)
			continue
		}

		s := string(lineBuf)
		s = strings.Trim(s, " ")
		fmt.Println(s)
		// write to channel
		channel <- []byte(s)
		b := websocket.NewFrame(true, websocket.OpcodeText, true, []byte(s))
		_, err = conn.Write(b)
		if err != nil {
			isExit <- true
		}

		select {
		case <-ctx.Done():
			return
		default:
			continue
		}
	}
}

func handleConnection(conn net.Conn) {
	ctx, cancel := context.WithCancel(context.Background())
	isExit := make(chan bool)
	channel := make(chan []byte)
	defer close(isExit)
	defer close(channel)

	go readFromTerminate(ctx, conn, channel, isExit)
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
