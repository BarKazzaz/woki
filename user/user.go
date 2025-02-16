package user

import (
	"fmt"
	"net"
	"strconv"
)

type User struct {
	Name       string
	Room       string
	Connection net.Conn
}

func (user *User) SendListRooms() {
	user.Connection.Write([]byte("getRooms\n"))
}

func (user *User) SendJoin(roomName string) {
	n := len(user.Name)
	msg := strconv.Itoa(n) + user.Name + "J" + roomName + "\n"
	user.Connection.Write([]byte(msg))
}

func (user *User) SendCreateRoom(roomName string) {
	n := len(user.Name)
	msg := strconv.Itoa(n) + user.Name + "C" + roomName + "\n"
	user.Connection.Write([]byte(msg))
}

func (user *User) Connect() {
	address, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	connection, err := net.DialTCP("tcp", nil, address)
	user.Connection = connection
	if err != nil {
		panic(err)
	}
}

func (user *User) SetUserName() {
	fmt.Printf("What is your name?\n")
	var userName string
	_, err := fmt.Scan(&userName)
	if err != nil {
		fmt.Printf("Error getting name: %v", err.Error())
		panic(err)
	}
	user.Name = userName
}

func (user *User) SendMessage(roomName, rest string) {
	msg := strconv.Itoa(len(user.Name))
	msg += user.Name + "M" + roomName + rest
	fmt.Printf("user.go:SendMessage(%v)\n", msg)
	user.Connection.Write([]byte(msg))
}

func (user *User) SendListUsers(roomName string) {
	// TODO: implement
	fmt.Printf("Cannot send list of users yet...\n")
}
