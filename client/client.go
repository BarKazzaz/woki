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

func handleCommand(cmd string) {
	command := strings.Split(cmd, " ")
	rest := strings.Join(command[1:], " ")
	fmt.Printf("command: %v, rest: %v\n", command[0], rest)
	switch command[0] {
	case "list":
		theUser.SendListRooms()
	case "join":
		theUser.SendJoin(rest)
	case "create":
		theUser.SendCreateRoom(rest)
	}
}

func connect() {
	address, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	connection, err := net.DialTCP("tcp", nil, address)
	theUser.Connection = connection
	if err != nil {
		panic(err)
	}
}

func setUserName() {
	fmt.Printf("What is your name?\n")
	var userName string
	_, err := fmt.Scan(&userName)
	if err != nil {
		fmt.Printf("Error getting name: %v", err.Error())
		panic(err)
	}
	theUser.Name = userName
}

func main() {
	connect()
	theUser.SendListRooms()
	buffer := make([]byte, 4096)
	roomsMsgLen, err := theUser.Connection.Read(buffer)
	if err != nil {
		fmt.Printf("Error getting rooms: %v", err.Error())
		panic(err)
	}

	setUserName()

	fmt.Printf("Avaliable rooms: %v\n", string(buffer[:roomsMsgLen]))

	fmt.Printf("Input room command: list/create/join\ne.g. join Main Loby\n")
	reader := bufio.NewReader(os.Stdin)
	cmd, err := reader.ReadString('\n')
	if err != nil {
		panic("Invalid input")
	}
	handleCommand(cmd)
	for {
		buffer = make([]byte, 4096)
		n, readErr := theUser.Connection.Read(buffer)
		timeoutTime := time.Now().Add(time.Minute * 5)
		if err := theUser.Connection.SetReadDeadline(timeoutTime); err != nil {
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
	defer theUser.Connection.Close()
}
