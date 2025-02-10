package chat

import (
	"errors"
	"net"
	"strings"
)

type User struct {
	Name       string
	Connection net.Conn
}

type ChatRoom struct {
	Chatters []User
}

func NewChatRoom() ChatRoom {
	c := ChatRoom{
		make([]User, 2),
	}
	return c
}

func (room ChatRoom) Join(user User) {
	room.Chatters = append(room.Chatters, user)
}

type Chat struct {
	Rooms map[string]ChatRoom
}

func (chat Chat) CreateRoom(roomName string) error {
	if chat.Rooms == nil {
		chat.Rooms = make(map[string]ChatRoom, 5)
	}
	if _, exists := chat.Rooms[roomName]; exists {
		return errors.New("Name already in use")
	}
	chat.Rooms[roomName] = NewChatRoom()
	return nil
}

func (chat Chat) JoinRoom(roomName string, user User) error {
	room, exists := chat.Rooms[roomName]
	if !exists {
		return errors.New("Room " + roomName + " doesn't exist")
	}
	room.Join(user)
	return nil
}

func (chat Chat) SendMessage(roomNameAndMessage []byte) error {
	roomNameAndMessageS := string(roomNameAndMessage)
	for roomName := range chat.Rooms {
		if strings.HasPrefix(roomNameAndMessageS, roomName) {
			msg := roomNameAndMessage[len(roomName):]
			chat.sendMessage(roomName, msg)
		}
	}
	return nil
}

func (chat Chat) sendMessage(roomName string, message []byte) error {
	room, exists := chat.Rooms[roomName]
	if !exists {
		return errors.New("Room " + roomName + " doesn't exist")
	}
	i := 0
	for i < len(room.Chatters) {
		room.Chatters[i].Connection.Write(message)
		i++
	}
	return nil
}
