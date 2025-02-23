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

func handleUserInput() {
	if len(theUser.Name) == 0 {
		theUser.SetUserName()
		return
	}
	reader := bufio.NewReader(os.Stdin)
	userInput, err := reader.ReadString('\n')
	if err != nil {
		panic("Invalid input")
	}
	userInput = strings.Split(userInput, "\n")[0]

	if userInput == "\\cmd help" {
		showInRoomCommands()
		return
	}

	insideARoom := currentRoom != ""
	split := strings.Split(userInput, " ")
	var command, rest string
	if !insideARoom {
		command = split[0]
		rest = strings.Join(split[1:], " ")
		handleOutOfRoomCommands(split[0], rest)
		return
	}

	command = "message"
	rest = userInput
	if strings.HasPrefix(userInput, "\\cmd ") {
		if len(split) < 2 {
			fmt.Printf("Invalid input: %v", userInput)
			return
		}
		command = split[1]
		if command == "room" {
			command = split[2]
			rest = strings.Join(split[3:], " ")
			handleOutOfRoomCommands(command, rest)
			return
		}
		rest = strings.Join(split[2:], " ")
	}
	handleInRoomCommands(command, rest)
}

func handleOutOfRoomCommands(command, rest string) {
	switch command {
	case "list":
		theUser.SendListRooms()
	case "join":
		theUser.SendJoin(rest)
	case "create":
		theUser.SendCreateRoom(rest)
	}
}

func handleInRoomCommands(command, rest string) {
	switch command {
	case "list":
		theUser.SendListUsers(currentRoom)
	case "change":
		theUser.SendJoin(rest)
	case "message":
		theUser.SendMessage(currentRoom, rest)
	}
}

func showInRoomCommands() {
	msg := "Type to send messages\n"
	msg += "\\cmd list\t\t\t(to list all users in the room)\n"
	msg += "\\cmd change\t\t\t(to change the room)\n"
	msg += "\\cmd room list\t\t\t(to list all rooms)\n"
	msg += "\\cmd room create\t\t(to create a new room)\n"
	msg += "\\cmd room join\t\t\t(to join a room - same as \\cmd change)\n"
	msg += "\\cmd help\t\t\t(to show this help message)\n"
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
	fmt.Printf("Welcome to Woki!\n")
	fmt.Printf("Type '\\cmd help' to see possible commands\n")
	theUser.Connect()
	defer theUser.Connection.Close()

	buffer := make([]byte, 4096)
	for {
		go handleUserInput()

		n, readErr := theUser.Connection.Read(buffer)
		setTimeout(readErr)
		msg := string(buffer[:n])
		if len(msg) > 1 {
			fmt.Printf("%v\n", msg)
		}

		prefix := "Error: "
		if strings.HasPrefix(msg, prefix) {
			continue
		}

		if strings.HasPrefix(msg, "Joined: ") {
			currentRoom = strings.Split(msg[:n], "Joined: ")[1]
			currentRoom = strings.Split(currentRoom, "\n")[0]
			fmt.Printf("You are in room: %v\n", currentRoom)
		}

		prefix = "User created: "
		if strings.HasPrefix(msg, prefix) {
			userName := strings.Split(msg[:n], prefix)[1]
			userName = strings.Split(userName, "\n")[0]
			theUser.Name = userName
			fmt.Printf("Changed name to: %v\n", userName)
		}
	}
}
