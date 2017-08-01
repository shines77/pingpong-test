//
// pingpong-test project client.go
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
	"time"
)

type FlagConfig struct {
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

var flagConfig FlagConfig

func init() {
	flag.IntVar(&flagConfig.processors, "thread-num", -1, "The number of work thread.")
	flag.StringVar(&flagConfig.protocol, "protocol", "tcp4", "The network type of host.")
	flag.StringVar(&flagConfig.host, "host", "localhost", "The IP address or domain name of host.")
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
	//fmt.Printf("tcp addr: %s\n", flagConfig.tcpAddr)
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

func setTCPSocketOptions(tcpConn *net.TCPConn, rdBufSize int, wrBufSize int, noDelay bool) {
	err := tcpConn.SetNoDelay(noDelay)
	if err != nil {
		fmt.Println("SetNoDelay() [nodelay=", flagConfig.nodelay, "] error: ", err)
	}

	if rdBufSize >= 0 {
		err = tcpConn.SetReadBuffer(rdBufSize)
		if err != nil {
			fmt.Println("SetReadBuffer() error: ", err)
		}
	}

	if wrBufSize >= 0 {
		err = tcpConn.SetWriteBuffer(wrBufSize)
		if err != nil {
			fmt.Println("SetWriteBuffer() error: ", err)
		}
	}
}

func setSocketOptions(conn net.Conn, rdBufSize int, wrBufSize int, noDelay bool) {
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

	setTCPSocketOptions(tcpConn, rdBufSize, rdBufSize, noDelay)
}

func ping(times int, pipeline int, lockChan chan bool) {
	tcpAddr, err := net.ResolveTCPAddr(flagConfig.protocol, flagConfig.tcpAddr)
	if err != nil {
		log.Fatal("get TCP error: ", err)
		panic(err)
	}
	tcpConn, err := net.DialTCP(flagConfig.protocol, nil, tcpAddr)
	if err != nil {
		log.Fatal("get DialTCP error: ", err)
		panic(err)
	}

	const READ_BUF_SIZE int = 160 * 1024
	const WRITE_BUF_SIZE int = 160 * 1024

	setTCPSocketOptions(tcpConn, READ_BUF_SIZE, WRITE_BUF_SIZE, flagConfig.nodelay)

	/*
		if tcpConn != nil {
			err := tcpConn.SetNoDelay(flagConfig.nodelay)
			if err != nil {
				log.Fatal("client tcp.SetNoDelay() [nodelay=", flagConfig.nodelay, "] error: ", err)
			}
		}
	*/

	if pipeline == 1 {
		for i := 0; i < times; i++ {
			nwrite, err := tcpConn.Write([]byte("Ping"))
			if err != nil {
				log.Fatal("get client write error: ", err)
			}
			var buff [4]byte
			if nwrite > 0 {
				nread, err := tcpConn.Read(buff[0:])
				if nread <= 0 {
					log.Fatal("Err: client read ", nread, " bytes")
				}
				if err != nil {
					log.Fatal("get client read error: ", err)
				}
			}
		}
		lockChan <- true
		tcpConn.Close()
	} else {
		for i := 0; i < times; i++ {
			for j := 0; j < pipeline; j++ {
				nwrite, err := tcpConn.Write([]byte("Ping"))
				if nwrite <= 0 {
					log.Fatal("Error: client write ", nwrite, " bytes")
				}
				if err != nil {
					log.Fatal("get client write error: ", err)
				}
			}
			var buff [4]byte
			for j := 0; j < pipeline; j++ {
				nread, err := tcpConn.Read(buff[0:])
				if nread <= 0 {
					log.Fatal("Error: client read ", nread, " bytes")
				}
				if err != nil {
					log.Fatal("get client read error: ", err)
				}
			}
		}
		lockChan <- true
		tcpConn.Close()
	}
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

	var pipeline int = flagConfig.pipeline
	// var totalPings int = 1000000
	var totalPings int = 500000
	var concurrentConnections int = 100
	var pingsPerConnection int = totalPings / (concurrentConnections * pipeline)
	if pingsPerConnection <= 0 {
		pingsPerConnection = 1
	}
	var actualTotalPings int = pingsPerConnection * concurrentConnections * pipeline
	lockChan := make(chan bool, concurrentConnections)

	tcpAddr, err := net.ResolveTCPAddr(flagConfig.protocol, flagConfig.tcpAddr)
	if err != nil {
		log.Fatal("get TCP error: ", err)
		panic(err)
	}
	fmt.Printf("Protocol: %s, Host: %s, Port: %d\n", tcpAddr.Network(), tcpAddr.IP.String(), tcpAddr.Port)
	fmt.Printf("\n")

	start := time.Now()
	for i := 0; i < concurrentConnections; i++ {
		go ping(pingsPerConnection, pipeline, lockChan)
	}
	for i := 0; i < concurrentConnections; i++ {
		<-lockChan
	}

	elapsed := time.Since(start).Seconds()
	fmt.Printf("actualTotalPings: %d, concurrentConnections: %d.\n", actualTotalPings, concurrentConnections)
	fmt.Printf("Elapsed time: %0.4f second(s), QPS: %0.1f query/sec.\n", elapsed, float64(actualTotalPings)/elapsed)
	fmt.Printf("\n")

	WaitEnter()
}
