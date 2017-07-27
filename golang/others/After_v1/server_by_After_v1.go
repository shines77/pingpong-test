//
// Author: After(QQ Group: 296561497)
// Date: 2016-08-15 15:58
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
	buf := make([]byte, 8192)
	reader := bufio.NewReaderSize(conn, 8192)
	writer := bufio.NewWriterSize(conn, 8192)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			fmt.Println("read error: ", err)
			return
		}

		for pos := 0; pos < n; pos += 64 {
			end := pos + 64
			if end > n {
				end = n
			}
			_, err1 := writer.Write(buf[pos:end])
			if err1 != nil {
				fmt.Println("write error: ", err1)
				return
			}
		}
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
