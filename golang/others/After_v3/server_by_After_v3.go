//
// Author: After(QQ Group: 296561497)
// Date: 2016-08-19 11:59
//

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReaderSize(conn, 2048)
	tmpChan := make(chan []byte, 20140)

	tcp := conn.(*net.TCPConn)
	if tcp != nil {
		tcp.SetNoDelay(true)
	}

	go func() {
		for {
			data := <-tmpChan
			remain := len(tmpChan)
			for i := 0; i < remain; i++ {
				left := <-tmpChan
				data = append(data, left...)
			}

			_, err1 := conn.Write(data[0:])
			if err1 != nil {
				fmt.Println("write error: ", err1)
				return
			}
		}

	}()

	buf := make([]byte, 8192)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			fmt.Println("read error: ", err)
			return
		}
		tmpChan <- buf[0:n]
	}
}

func main() {
	ln, err := net.Listen("tcp", ":5178")
	if err != nil {
		panic(err)
	}

	fmt.Println("listen 5178 ok")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("get client connection error: ", err)
		}

		go handleConnection(conn)
	}
}
