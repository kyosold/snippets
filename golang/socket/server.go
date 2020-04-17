package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"strings"
	"syscall"
	"time"
)

const readTimeout int = 20
const writeTimeout int = 20

const MAXBUFSIZE = 4096

func main() {
	flag.Parse()
	if flag.NArg() != 2 {
		panic("Usage: host port")
	}
	hostAndPort := fmt.Sprintf("%s:%s", flag.Arg(0), flag.Arg(1))
	listener := initServer(hostAndPort)
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic("Accept: " + err.Error())
		}
		go connectionHandler(conn)
	}
}

func initServer(hostAndPort string) *net.TCPListener {
	serverAddr, err := net.ResolveTCPAddr("tcp", hostAndPort)
	if err != nil {
		panic("ERROR: Resolving address:port failed: " + hostAndPort + ", " + err.Error())
	}
	listener, err := net.ListenTCP("tcp", serverAddr)
	if err != nil {
		panic("ERROR: ListenTCP: " + err.Error())
	}
	fmt.Println("Listening to: ", listener.Addr().String())
	return listener
}

func connectionHandler(conn net.Conn) {
	connFrom := conn.RemoteAddr().String()
	fmt.Println("Connect from: ", connFrom)
	sayGreeting(conn)
	for {
		ibuf := make([]byte, MAXBUFSIZE+1)

		// 设置读超时
		conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(readTimeout)))
		nr, err := conn.Read(ibuf[0:MAXBUFSIZE])
		ibuf[MAXBUFSIZE] = 0 // to prevent overflow
		if err == nil {
			showMsg("r", nr, err, ibuf)
		} else if err == syscall.EAGAIN { // try again
			continue
		} else if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			fmt.Println("[ERR] Read Timeout")
			continue
		} else if err == io.EOF {
			fmt.Println("Client CLOSED")
			goto DISCONNECT
		} else {
			goto DISCONNECT
		}

		// 设置写超时
		conn.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(writeTimeout)))
		outStr := "250 OK\r\n"
		_, err = conn.Write([]byte(outStr))
		if err == nil {
			outStr = strings.Trim(outStr, "\r\n")
			showMsg("w", len(outStr), err, []byte(outStr))
		} else if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			fmt.Println("[ERR] Write Timeout")
			continue
		} else {
			goto DISCONNECT
		}
	}

DISCONNECT:
	err := conn.Close()
	fmt.Println("Closed connection: ", connFrom)
	if err != nil {
		fmt.Println("Close: ", err.Error())
	}
}

func sayGreeting(to net.Conn) {
	obufstr := "Let's GO!\n"
	wrote, err := to.Write([]byte(obufstr))
	if err != nil {
		panic("Write: wrote " + string(wrote) + " bytes.")
	}
}

func showMsg(t string, len int, err error, msg []byte) {
	if len > 0 {
		if t == "r" {
			fmt.Print("[READ] ")
		} else {
			fmt.Print("[WRITE] ")
		}
		fmt.Print("<", len, ":")
		for i := 0; ; i++ {
			if i >= len || msg[i] == 0 {
				break
			}
			fmt.Printf("%c", msg[i])
		}
		fmt.Print(">\n")
	}
}
