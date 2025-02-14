package user

import (
	"net"
	"strconv"
)

type User struct {
	Name       string
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
