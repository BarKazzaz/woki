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
	case 'S':
		nameExists := theChat.HasUser(name)
		if nameExists {
			conn.Write([]byte("Error: name is already taken\n"))
			return
		}
		conn.Write([]byte("User created: " + name + "\n"))
		var defaultRoom string
		for name := range theChat.Rooms {
			defaultRoom = name
			break
		}
		err := theChat.JoinRoom(defaultRoom, &user)
		if err != nil {
			conn.Write([]byte("Error:" + err.Error() + "\n"))
			return
		}
		conn.Write([]byte("Joined: " + defaultRoom + "\n"))
	case 'C':
		err := theChat.CreateRoom(arg)
		if err != nil {
			conn.Write([]byte("Error:" + err.Error() + "\n"))
			return
		}
		conn.Write([]byte("Room created: " + arg + "\n"))
	case 'J':
		err := theChat.JoinRoom(arg, &user)
		if err != nil {
			conn.Write([]byte("Error:" + err.Error() + "\n"))
			return
		}
		conn.Write([]byte("Joined: " + arg + "\n"))
	case 'M':
		err := theChat.SendMessage(&user, arg)
		if err != nil {
			conn.Write([]byte("Error:" + err.Error() + "\n"))
			return
		}
		conn.Write([]byte("\n"))
	case 'R':
		rooms := theChat.GetRoomsList()
		s := "[" + strings.Join(rooms, ",") + "]\n"
		conn.Write([]byte(s))
	case 'U':
		var users []string
		for _, u := range theChat.Rooms[arg].Chatters {
			if u == nil {
				continue
			}
			users = append(users, u.Name)
		}
		s := "[" + strings.Join(users, ",") + "]\n"
		conn.Write([]byte(s))
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
			theChat.RemoveUser(connection)
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
	mainRoomName := "led"
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
