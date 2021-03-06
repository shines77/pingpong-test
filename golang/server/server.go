//
// pingpong-test project server.go
//
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type sFlagConfig struct {
	processors int
	protocol   string
	host       string
	port       int
	tcpAddr    string
	pipeline   int
	nodelayStr string
	nodelay    bool
	args       []string
}

var flagConfig sFlagConfig

func init() {
	flag.IntVar(&flagConfig.processors, "thread-num", -1, "The number of work thread.")
	flag.StringVar(&flagConfig.protocol, "protocol", "tcp4", "The network type of host.")
	flag.StringVar(&flagConfig.host, "host", "", "The IP address or domain name of host.")
	flag.IntVar(&flagConfig.port, "port", 5178, "The port of host.")
	flag.IntVar(&flagConfig.pipeline, "pipeline", 1, "The pipeline of ping one time.")
	flag.StringVar(&flagConfig.nodelayStr, "nodelay", "false", "TCP is setting nodelay mode? options is [0,1] or [true,false].")
}

func initArgs() {
	if strings.TrimSpace(flagConfig.host) == "" {
		flagConfig.tcpAddr = fmt.Sprintf(":%d", flagConfig.port)
	} else {
		flagConfig.tcpAddr = fmt.Sprintf("%s:%d", flagConfig.host, flagConfig.port)
	}
	if flagConfig.pipeline <= 0 {
		flagConfig.pipeline = 1
	}
	flagConfig.nodelay = parseBool(&flagConfig.nodelayStr, false)
	flagConfig.args = flag.Args()
}

func printArgs() {
	fmt.Printf("The input arguments:\n\n")
	fmt.Printf("thread-num: %d\n", flagConfig.processors)
	fmt.Printf("protocol: %s\n", formatString(&flagConfig.protocol))
	fmt.Printf("host: %s\n", formatString(&flagConfig.host))
	fmt.Printf("port: %d\n", flagConfig.port)
	fmt.Printf("pipeline: %d\n", flagConfig.pipeline)
	fmt.Printf("nodelay: %s\n", strconv.FormatBool(flagConfig.nodelay))
	fmt.Printf("other args: %s\n", flagConfig.args)
	fmt.Printf("tcp addr: %s\n", flagConfig.tcpAddr)
	fmt.Printf("\n")
}

func parseArgs() {
	flag.Parse()
	initArgs()
	printArgs()

	if flagConfig.port <= 0 && flagConfig.port > 65535 {
		err := fmt.Errorf("The port out of range [1, 65535]: %d\n", flagConfig.port)
		fmt.Println(err)
		os.Exit(1)
	}
}

func setTCPSocketOptions(tcpConn *net.TCPConn, readBufSize int, writeBufSize int, noDelay bool) {
	err := tcpConn.SetNoDelay(noDelay)
	if err != nil {
		fmt.Println("SetNoDelay() [nodelay=", flagConfig.nodelay, "] error: ", err)
	}

	if readBufSize >= 0 {
		err = tcpConn.SetReadBuffer(readBufSize)
		if err != nil {
			fmt.Println("SetReadBuffer() error: ", err)
		}
	}

	if writeBufSize >= 0 {
		err = tcpConn.SetWriteBuffer(writeBufSize)
		if err != nil {
			fmt.Println("SetWriteBuffer() error: ", err)
		}
	}
}

func setSocketOptions(conn net.Conn, readBufSize int, writeBufSize int, noDelay bool) {
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

	setTCPSocketOptions(tcpConn, readBufSize, writeBufSize, noDelay)
}

func pong(conn net.Conn, pipeline int) {
	var buf [4]byte

	if pipeline == 1 {
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
	} else {
		for {
			for j := 0; j < pipeline; j++ {
				_, err := conn.Read(buf[0:])
				if err != nil {
					return
				}
			}

			for j := 0; j < pipeline; j++ {
				_, err := conn.Write([]byte("Pong"))
				if err != nil {
					return
				}
			}
		}
	}
}

func handleClient(conn net.Conn, pipeline int) {
	defer conn.Close()

	const readBufSize int = 160 * 1024
	const writeBufSize int = 160 * 1024

	setSocketOptions(conn, readBufSize, writeBufSize, flagConfig.nodelay)

	/*
		tcp := conn.(*net.TCPConn)
		if tcp != nil {
			err := tcp.SetNoDelay(flagConfig.nodelay)
			if err != nil {
				log.Fatal("client tcp.SetNoDelay() [nodelay=", flagConfig.nodelay, "] error: ", err)
			}
		}
	*/

	pong(conn, pipeline)
}

func main() {
	log.SetOutput(os.Stdout)

	parseArgs()

	fullProcessors := runtime.GOMAXPROCS(flagConfig.processors)
	if flagConfig.processors <= 0 {
		fmt.Printf("Processors: %d / %d\n", fullProcessors, fullProcessors)
	} else {
		fmt.Printf("Processors: %d / %d\n", flagConfig.processors, fullProcessors)
	}
	fmt.Printf("\n")

	tcpAddr, err := net.ResolveTCPAddr(flagConfig.protocol, flagConfig.tcpAddr)
	if err != nil {
		log.Fatal("get TCP error: ", err)
		panic(err)
	}
	listener, err := net.ListenTCP(flagConfig.protocol, tcpAddr)
	if err != nil {
		log.Fatal("get listen error: ", err)
		panic(err)
	}

	fmt.Printf("Protocol: %s, Host: %s, Port: %d\n", tcpAddr.Network(), tcpAddr.IP.String(), tcpAddr.Port)
	fmt.Printf("\n")
	fmt.Printf("Listening %s %s ...", listener.Addr().Network(), listener.Addr().String())

	pipeline := flagConfig.pipeline
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("get client connection error: ", err)
		}
		go handleClient(conn, pipeline)
	}
}
