//
// Author: TW2(QQ Group: 296561497)
// Date: 2016-08-16 17:02
//

package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	//fmt.Println("new")

	buf := make([]byte, 8192)
	timeout := 30

	for {
		conn.SetReadDeadline(time.Now().Add(time.Millisecond * time.Duration(timeout)))
		n, err := conn.Read(buf)
		conn.SetReadDeadline(time.Now())
		if err != nil {
			//fmt.Println("read error: ", err)
			e, ok := err.(net.Error)
			if !ok || !e.Temporary() {
				fmt.Println("read error: ", err, " return")
				return
			}
		}

		for pos := 0; pos < n; { //pos += 64
			end := pos + 64
			if end > n {
				end = n
			}
			//conn.SetWriteDeadline(time.Now().Add(time.Millisecond * time.Duration(timeout)))
			sended, err := conn.Write(buf[pos:end])
			//conn.SetWriteDeadline(time.Now())
			if err != nil {
				//fmt.Println("write error: ", err)

				e, ok := err.(net.Error)
				if !ok || !e.Temporary() {
					fmt.Println("write error: ", err, " return")
					return
				}
			}

			pos += sended
		}
	}
	return
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
