package main

import(
	"fmt"
	"log"
	"net"
)

func setSocketOptions(conn net.Conn, bufSize int, noDelay bool) {
	//
	// See:	http://tonybai.com/2015/11/17/tcp-programming-in-golang/
	// See: https://golang.org/pkg/net/#TCPConn.SetNoDelay
	//
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		// error handle
		fmt.Println("TCPConn type assertion error.")
		return
	}

	err := tcpConn.SetNoDelay(noDelay)
	if err != nil {
		fmt.Println("SetNoDelay() error: ", err)
	}

	if bufSize < 0 {
		return
	}

	err = tcpConn.SetReadBuffer(bufSize)
	if err != nil {
		fmt.Println("SetReadBuffer() error: ", err)
	}

	err = tcpConn.SetWriteBuffer(bufSize)
	if err != nil {
		fmt.Println("SetWriteBuffer() error: ", err)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	BUF_SIZE := 160 * 1024

	setSocketOptions(conn, 160 * 1024, false);
	setSocketOptions(conn, -1, false);

	buf := make([]byte, BUF_SIZE)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error: ", err)
			return
		}

		if n != BUF_SIZE {
			fmt.Println("Receive Bytes: ", n)
		}

//		_, err1 := conn.Write(buf[0:n])
//		if err1 != nil {
//			fmt.Println("write error: ", err)
//			return
//		}

		for pos := 0; pos < n; pos += 64 {
			end := pos + 64
			if end > n {
				end = n
			}
			_, err1 := conn.Write(buf[pos:end])
			if err1 != nil {
				fmt.Println("write error: ", err)
				return
			}
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
