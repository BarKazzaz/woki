package main

import (
	"fmt"
	"net"
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
	buffer := make([]byte, 4094)
	n, err := connection.Read(buffer)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Read %d bytes\nMessage: %v\n", n, string(buffer[:n]))

	connection.Write([]byte("3barJyoyo\n"))
	buffer = make([]byte, 4094)
	n, err = connection.Read(buffer)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Read %d bytes\nMessage: %v\n", n, string(buffer[:n]))
	connection.Write([]byte("4done\n"))
	connection.Close()
}
