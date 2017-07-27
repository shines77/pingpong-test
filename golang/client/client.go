//
// pingpong-test project client.go
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
	"time"
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
	__hostIP = flag.String("host", "localhost", "The IP address or domain name of host.")
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
	//fmt.Printf("tcp addr: %s\n", __tcpAddr)
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

func ping(times int, lockChan chan bool) {
	tcpAddr, err := net.ResolveTCPAddr(*__protocol, __tcpAddr)
	if err != nil {
		panic(err)
	}
	conn, err := net.DialTCP(*__protocol, nil, tcpAddr)
	if err != nil {
		panic(err)
	}

	if conn != nil {
		err := conn.SetNoDelay(__nodelay)
		if err != nil {
			log.Fatal("client tcp.SetNoDelay(", __nodelay, ") error: ", err)
		}
	}

	for i := 0; i < times; i++ {
		write_bytes, err := conn.Write([]byte("Ping"))
		if err != nil {
			log.Fatal("get client write error: ", err)
		}
		var buff [4]byte
		if write_bytes > 0 {
			read_bytes, err := conn.Read(buff[0:])
			if read_bytes <= 0 {
				log.Fatal("Err: client read ", read_bytes, " bytes")
			}
			if err != nil {
				log.Fatal("get client read error: ", err)
			}
		}
	}
	lockChan <- true
	conn.Close()
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

	// var totalPings int = 1000000
	var totalPings int = 300000
	var concurrentConnections int = 100
	var pingsPerConnection int = totalPings / concurrentConnections
	var actualTotalPings int = pingsPerConnection * concurrentConnections
	lockChan := make(chan bool, concurrentConnections)

	tcpAddr, err := net.ResolveTCPAddr(*__protocol, __tcpAddr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Protocol: %s, Host: %s, Port: %d\n", tcpAddr.Network(), tcpAddr.IP.String(), tcpAddr.Port)
	fmt.Printf("\n")

	start := time.Now()
	for i := 0; i < concurrentConnections; i++ {
		go ping(pingsPerConnection, lockChan)
	}
	for i := 0; i < concurrentConnections; i++ {
		<-lockChan
	}

	elapsed := time.Since(start).Seconds()
	fmt.Printf("actualTotalPings: %d, concurrentConnections: %d.\n", actualTotalPings, concurrentConnections)
	fmt.Printf("Elapsed time: %0.4f second(s), QPS: %0.1f query/sec.\n", elapsed, float64(actualTotalPings)/elapsed)
	fmt.Printf("\n")

	waitEnter()
}