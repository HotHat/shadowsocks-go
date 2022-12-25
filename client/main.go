package main

import (
	"context"
	"fmt"
	"net"
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
	for {
		n, err := conn.Read(buf)
		if err != nil {
			isExit <- true
			break
		}
		fmt.Println("Read:", n, " data:", buf[:n])
		channel <- buf[:n]
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
