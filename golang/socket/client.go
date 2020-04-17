package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const dialTimeout int = 2
const readTimeout int = 10
const writeTimeout int = 10

func main() {
	flag.Parse()
	if flag.NArg() != 2 {
		panic("Usage: host port")
	}
	hostAndPort := fmt.Sprintf("%s:%s", flag.Arg(0), flag.Arg(1))

	// 打开连接
	conn, err := net.DialTimeout("tcp", hostAndPort, time.Duration(dialTimeout)*time.Second)
	if err != nil {
		// 由于目标计算机积极拒绝而无法创建连接
		fmt.Println("Error dialing", err.Error())
		return
	}
	defer conn.Close()
	fmt.Println("Connect OK")

	inReader := bufio.NewReader(os.Stdin)
	fmt.Println("First, what's your name ?")
	clientName, _ := inReader.ReadString('\n')
	trimmedClient := strings.Trim(clientName, "\r\n") // Windows下用"\r\n", Linux下用"\n"

	// 设置从终端读取内容，用于在主线程设置超时。
	inStdinChan := make(chan string)
	go func() {
		inReader := bufio.NewReader(os.Stdin)
		for {
			fmt.Println("What to send to the server ? Type Q to Quit.")
			o, _ := inReader.ReadString('\n')
			if err != nil {
				fmt.Printf("inStdinChan error: %v", err)
				return
			}
			inStdinChan <- o
		}
	}()

	// 给服务器发送信息直到退出
	for {
		var input string
		var trimmedInput string
		select {
		case input = <-inStdinChan:
			trimmedInput = strings.Trim(input, "\r\n")
			if trimmedInput == "Q" {
				return
			}
		case <-time.After(time.Duration(readTimeout) * time.Second):
			fmt.Printf("[Err]: Read Timeout(%d Seconds)", readTimeout)
			return
		}

		// 设置服务器写超时
		conn.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(writeTimeout)))
		str := trimmedClient + " says: " + trimmedInput
		n, err := conn.Write([]byte(str))
		if err != nil {
			fmt.Println("> WRITE FAIL:", err)
		}
		fmt.Printf("> WRITE(%d): %s\n", n, str)

		// 设置服务器读超时
		conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(readTimeout)))
		ibuf := make([]byte, 4096)
		n, err = conn.Read(ibuf)
		if err != nil {
			fmt.Println("< READ FAIL:", err)
		}
		fmt.Printf("< READ(%d); %s", n, string(ibuf))
	}
}
