package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	target := "localhost:9999"

	raddr, err := net.ResolveTCPAddr("tcp", target)
	if err != nil {
		log.Fatal(err)
	}

	// 和服务端建立连接
	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		log.Fatal(err)
	}

	//conn.SetNoDelay(false) // 如果打开这行代码，则禁用TCP_NODELAY，打开Nagle算法

	fmt.Println("Sending Gophers down the pipe...")

	for i := 0; i < 5; i++ {
		// 发送数据
		_, err = conn.Write([]byte("GOPHER"))
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(1000)
	}
}
