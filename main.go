package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"woki/chat"
)

var theChat = chat.Chat{}

func parse(requestBody []byte, conn net.Conn) {
	// TODO: support names > 9 ([0] means only 1 char unless we go the hexa way of numbering)
	var nameLength int = int(requestBody[0] - byte('0'))
	if nameLength < 3 {
		return
	}
	var nameLastIndex int = nameLength + 1

	name := string(requestBody[1:nameLastIndex])
	command := requestBody[nameLastIndex] // Create | Join | Message
	trailingNewLineIndex := 0
	for trailingNewLineIndex < len(requestBody) {
		if requestBody[trailingNewLineIndex] == '\n' {
			break
		}
		trailingNewLineIndex++
	}

	arg := string(requestBody[nameLastIndex+1 : trailingNewLineIndex]) // roomName | message
	fmt.Printf("%v name: %v command: %v arg: %v\n", strconv.Itoa(nameLength), name, string(command), arg)

	user := chat.User{}
	user.Name = string(name)
	user.Connection = conn

	switch command {
	case 'C':
		conn.Write([]byte("Creating room " + arg + "\n"))
		err := theChat.CreateRoom(arg)
		if err != nil {
			conn.Write([]byte("Error:" + err.Error() + "\n"))
		}
	case 'J':
		fmt.Printf("Trying to join " + name + " to room: " + arg)
		conn.Write([]byte("Joining room " + arg + "\n"))
		err := theChat.JoinRoom(arg, user)
		if err != nil {
			conn.Write([]byte("Error:" + err.Error() + "\n"))
		}
	case 'M':
		err := theChat.SendMessage(arg)
		if err != nil {
			conn.Write([]byte("Error:" + err.Error() + "\n"))
		}
	default:
		conn.Write([]byte("Unknown command:" + string(command) + "\n"))
	}
}

func hey(c net.Conn) {
	c.Write([]byte("Hey"))
}

func handleConnetion(connection net.Conn) {
	defer connection.Close()
	for {
		buff := make([]byte, 4096)
		i, err := connection.Read(buff)
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			fmt.Println("EOF")
			break
		}
		fmt.Println("Read len bytes:", i)
		parse(buff, connection)
	}
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Println("Serving", listener.Addr().Network(), listener.Addr().String())
	for {
		connection, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnetion(connection)
	}
}
