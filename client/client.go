package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	address, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	connection, err := net.DialTCP("tcp", nil, address)
	if err != nil {
		panic(err)
	}
	connection.Write([]byte("3barCyoyo\n"))
	buffer := make([]byte, 4096)
	n, readErr := connection.Read(buffer)
	if readErr != nil {
		panic(readErr)
	}
	fmt.Printf("Read %d bytes\nMessage: %v\n", n, string(buffer[:n]))

	connection.Write([]byte("3barJyoyo\n"))
	connection.Write([]byte(io.EOF.Error()))
	for {
		buffer = make([]byte, 4096)
		n, readErr = connection.Read(buffer)
		timeoutTime := time.Now().Add(time.Second * 5)
		if err := connection.SetReadDeadline(timeoutTime); err != nil {
			fmt.Printf("Timeout")
			break
		}
		if readErr != nil {
			if netErr, ok := readErr.(net.Error); ok && netErr.Timeout() {
				fmt.Printf("Timeout\n")
				break
			} else {
				panic(readErr)
			}
		}
		msg := string(buffer[:n])
		fmt.Printf("Read %d bytes\nMessage: %v\n", n, msg)
	}
	defer connection.Close()
}
