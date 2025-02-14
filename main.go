package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"woki/chat"
	"woki/user"
)

var theChat = chat.Chat{}

func parse(requestBody []byte, conn net.Conn) {
	trailingNewLineIndex := 0
	for trailingNewLineIndex < len(requestBody) {
		if requestBody[trailingNewLineIndex] == '\n' {
			break
		}
		trailingNewLineIndex++
	}
	line := string(requestBody[:trailingNewLineIndex])
	fmt.Printf("line: %v (%v)\n", line, len(line))
	if line == "getRooms" {
		rooms := theChat.GetRoomsList()
		fmt.Printf("getRooms: %v\n", rooms)
		s := "[" + strings.Join(rooms, ",") + "]\n"
		conn.Write([]byte(s))
		return
	}
	// TODO: support names > 9 ([0] means only 1 char unless we go the hexa way of numbering)
	var nameLength int = int(requestBody[0] - byte('0'))
	if nameLength < 3 {
		return
	}
	var nameLastIndex int = nameLength + 1

	name := string(requestBody[1:nameLastIndex])
	command := requestBody[nameLastIndex] // Create | Join | Message

	arg := string(requestBody[nameLastIndex+1 : trailingNewLineIndex]) // roomName | message
	fmt.Printf("%v name: %v command: %v arg: %v\n", strconv.Itoa(nameLength), name, string(command), arg)

	user := user.User{}
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
		fmt.Printf("Trying to join " + name + " to room: " + arg + "\n")
		conn.Write([]byte("Joining room " + arg + "\n"))
		err := theChat.JoinRoom(arg, user)
		if err != nil {
			conn.Write([]byte("Error:" + err.Error() + "\n"))
		}
	case 'M':
		err := theChat.SendMessage(user, arg)
		if err != nil {
			conn.Write([]byte("Error:" + err.Error() + "\n"))
		}
	default:
		conn.Write([]byte("Unknown command:" + string(command) + "\n"))
	}
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
	mainRoomName := "Main Lobby"
	if !theChat.HasRoom(mainRoomName) {
		theChat.CreateRoom(mainRoomName)
	}
	for {
		connection, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnetion(connection)
	}
}
