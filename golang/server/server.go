//
// pingpong-test project server.go
//
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"runtime"
	"strconv"
	"strings"
)

var __processors *int
var __protocol *string
var __hostIP *string
var __hostPort *int
var __nodelay_str *string
var __nodelay bool
var __tcpAddr string

func init() {
	__processors = flag.Int("thread-num", -1, "The number of work thread.")
	__protocol = flag.String("protocol", "tcp4", "The network type of host.")
	__hostIP = flag.String("host", "", "The IP address or domain name of host.")
	__hostPort = flag.Int("port", 5178, "The port of host.")
	__nodelay_str = flag.String("nodelay", "true", "Whether TCP use nodelay mode, options is [0,1] or [true,false].")
}

func init_args() {
	if strings.TrimSpace(*__hostIP) == "" {
		__tcpAddr = fmt.Sprintf(":%d", *__hostPort)
	} else {
		__tcpAddr = fmt.Sprintf("%s:%d", *__hostIP, *__hostPort)
	}
	__nodelay = parseBool(__nodelay_str, false)
}

func print_args() {
	fmt.Printf("The input arguments:\n\n")
	fmt.Printf("thread-num: %d\n", *__processors)
	fmt.Printf("protocol: %s\n", formatString(__protocol))
	fmt.Printf("host: %s\n", formatString(__hostIP))
	fmt.Printf("port: %d\n", *__hostPort)
	fmt.Printf("nodelay: %s\n", strconv.FormatBool(__nodelay))
	fmt.Printf("other args: %s\n", flag.Args())
	fmt.Printf("tcp addr: %s\n", __tcpAddr)
	fmt.Printf("\n")
}

func parse_args() {
	flag.Parse()
	init_args()
	print_args()

	if *__hostPort <= 0 && *__hostPort > 65535 {
		fmt.Errorf("The port out of range [1, 65535]: %d\n", *__hostPort)
		return
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	tcp := conn.(*net.TCPConn)
	if tcp != nil {
		err := tcp.SetNoDelay(__nodelay)
		if err != nil {
			log.Fatal("client tcp.SetNoDelay(", __nodelay, ") error: ", err)
		}
	}

	var buf [4]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		if n > 0 {
			_, err = conn.Write([]byte("Pong"))
			if err != nil {
				return
			}
		}
	}
}

func main() {
	parse_args()

	fullProcessors := runtime.GOMAXPROCS(*__processors)
	if *__processors <= 0 {
		fmt.Printf("Processors: %d / %d\n", fullProcessors, fullProcessors)
	} else {
		fmt.Printf("Processors: %d / %d\n", *__processors, fullProcessors)
	}
	fmt.Printf("\n")

	tcpAddr, err := net.ResolveTCPAddr(*__protocol, __tcpAddr)
	if err != nil {
		log.Fatal("get TCP error: ", err)
		panic(err)
	}
	listener, err := net.ListenTCP(*__protocol, tcpAddr)
	if err != nil {
		log.Fatal("get listen error: ", err)
		panic(err)
	}

	fmt.Printf("Protocol: %s, Host: %s, Port: %d\n", tcpAddr.Network(), tcpAddr.IP.String(), tcpAddr.Port)
	fmt.Printf("\n")
	fmt.Printf("Listening %s %s ...", listener.Addr().Network(), listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("get client connection error: ", err)
		}
		go handleClient(conn)
	}
}
