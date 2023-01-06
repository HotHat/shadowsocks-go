package main

import (
	"fmt"
	"log"
	"net"
)

const BufLen = 13

func main() {
	port := ":" + "9999"

	// 创建监听
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		// 接收新的连接
		c, err := l.Accept()
		if err != nil {
			log.Println(err)
			return
		}

		// 处理新的连接
		go handleConnection(c)
	}
}
func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	buf := make([]byte, BufLen)
	readLn := 0
	times := 1
	for {
		fmt.Println("loop times:", times)
		times += 1

		t := buf[readLn:]
		// 读取数据
		netData, err := c.Read(t)
		readLn += netData
		fmt.Println("current len:", netData)
		fmt.Println("read len:", readLn)
		fmt.Println("current read:", string(t[:netData]))
		fmt.Println("session read:", string(buf[:readLn]))

		if err != nil {
			log.Println(err)
			return
		}

		if readLn < BufLen {
			fmt.Printf("not fixed %d continue read\n", BufLen)
			continue
		}

		readLn = 0

		//cdata := strings.TrimSpace(netData)
		//if cdata == "GOPHER" {
		//	c.Write([]byte("GopherAcademy Advent 2023!"))
		//}
		//
		//if cdata == "EXIT" {
		//	break
		//}
	}
	c.Close()
}
