package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
	"woki/user"
)

var theUser = user.User{}
var currentRoom string

func handleCommands(cmd string) {
	insideARoom := currentRoom != ""
	split := strings.Split(cmd, " ")
	var command, rest string
	if insideARoom {
		if strings.HasPrefix(cmd, "\\cmd ") {
			if len(split) < 2 {
				fmt.Printf("Invalid input: %v", cmd)
				return
			}
			command = split[1]
			rest = strings.Join(split[2:], " ")
		} else {
			command = "message"
			rest = cmd
		}
		handleInRoomCommands(command, rest)
		return
	}
	command = split[0]
	rest = strings.Join(split[1:], " ")
	handleOutOfRoomCommands(split[0], rest)
}

func handleOutOfRoomCommands(command, rest string) {
	fmt.Printf("command: %v, rest: %v\n", command, rest)
	switch command {
	case "list":
		theUser.SendListRooms()
	case "join":
		currentRoom = rest
		theUser.SendJoin(currentRoom)
	case "create":
		theUser.SendCreateRoom(rest)
	}
}

func handleInRoomCommands(command, rest string) {
	fmt.Printf("command: %v, rest: %v\n", command, rest)
	switch command {
	case "list":
		fmt.Printf("asking for list of users...\n")
		theUser.SendListUsers(currentRoom)
	case "change":
		theUser.SendJoin(rest)
	case "message":
		theUser.SendMessage(currentRoom, rest)
	}
}

func showOutOfRoomCommands() {
	fmt.Printf("Commands: list/create/join\ne.g. join Main Loby\nOR\ncreate Led Zepplin Fans\n")
}

func showInRoomCommands() {
	msg := "Type to send messages\n"
	msg += "Commands:					change,list\n"
	msg += "Hello							(to send 'Hello' to other chatters)\n"
	msg += "\\cmd list				(to list all users in the room)\n"
	msg += "\\cmd change			(to change the room)\n"
	fmt.Printf(msg)
}

func setTimeout(readErr error) {
	timeoutTime := time.Now().Add(time.Minute * 5)
	if err := theUser.Connection.SetReadDeadline(timeoutTime); err != nil {
		fmt.Printf("Timeout..\n")
	}
	if readErr != nil {
		if netErr, ok := readErr.(net.Error); ok && netErr.Timeout() {
			fmt.Printf("Timeout\n")
		} else {
			panic(readErr)
		}
	}
}

func main() {
	theUser.Connect()
	defer theUser.Connection.Close()
	buffer := make([]byte, 4096)

	theUser.SetUserName()

	theUser.SendListRooms()
	roomsMsgLen, err := theUser.Connection.Read(buffer)
	if err != nil {
		fmt.Printf("Error getting rooms: %v", err.Error())
		panic(err)
	}

	fmt.Printf("Welcome to Woki!\nAvaliable rooms: %v\n", string(buffer[:roomsMsgLen]))

	for {
		buffer = make([]byte, 4096) // TODO: don't create in every iteration
		if currentRoom != "" {
			fmt.Printf("You are in room: %v", currentRoom)
			showInRoomCommands()
		} else {
			showOutOfRoomCommands()
		}
		reader := bufio.NewReader(os.Stdin)
		cmd, err := reader.ReadString('\n')
		if err != nil {
			panic("Invalid input")
		}
		handleCommands(cmd)
		n, readErr := theUser.Connection.Read(buffer)
		setTimeout(readErr)
		msg := string(buffer[:n])
		fmt.Printf("%v\n", msg)

		// TODO: move things to handleResponse()
		if strings.HasPrefix(msg, "Joined: ") {
			currentRoom = strings.Split(msg[:n], "Joined: ")[1]
			fmt.Printf("Changed current room to: %v", currentRoom)
		}
	}
}
