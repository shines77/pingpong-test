package main

import(
	"fmt"
	"log"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 8192)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error: ", err)
			return
		}
		_, err1 := conn.Write(buf[0:n])
		if err1 != nil {
			fmt.Println("write error: ", err)
			return
		}
	}
}

func main() {
	ln, err := net.Listen("tcp", ":9003")
	if err != nil {
		panic(err)
	}

	fmt.Println("listen 9003 ok")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("get client connection error: ", err)
		}

		go handleConnection(conn)
	}
}
