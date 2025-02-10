package main

import (
	"fmt"
	"io"
	"net"
	"woki/chat"
)

var theChat = chat.Chat{}

func parse(requestBody []byte, conn net.Conn) {
	// TODO: support names > 9 ([0] means only 1 char)
	nameLength := int(requestBody[0]) - int(byte('0')) + 1

	name := requestBody[1:nameLength]
	command := requestBody[nameLength] // Create | Join | Message
	arg := requestBody[nameLength+1:]  // roomName | message
	fmt.Printf("%v name %v, command %v arg %v", string(nameLength), string(name), string(command), string(arg))

	user := chat.User{}
	user.Name = string(name)
	user.Connection = conn

	switch command {
	case 'C':
		room := string(arg)
		conn.Write([]byte("Creating room " + room))
		theChat.CreateRoom(room)
	case 'J':
		room := string(arg)
		conn.Write([]byte("Joining room " + room))
		theChat.JoinRoom(room, user)
	case 'M':
		theChat.SendMessage(arg)
	default:
		conn.Write([]byte("Unknown command:" + string(command)))
	}
}

func hey(c net.Conn) {
	c.Write([]byte("Hey"))
}

func handleConnetion(connection net.Conn) {
	buff := make([]byte, 4096)
	defer connection.Close()
	for {
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
