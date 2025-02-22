package chat

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"woki/user"
)

type ChatRoom struct {
	Chatters []*user.User
}

func NewChatRoom() ChatRoom {
	c := ChatRoom{
		make([]*user.User, 0),
	}
	return c
}

func (room *ChatRoom) Join(user *user.User) {
	room.Chatters = append(room.Chatters, user)
	fmt.Printf("%v", room)
}

type Chat struct {
	Rooms map[string]*ChatRoom
}

func (chat *Chat) HasRoom(roomName string) bool {
	_, roomExists := chat.Rooms[roomName]
	return roomExists
}

func (chat *Chat) CreateRoom(roomName string) error {
	if chat.Rooms == nil {
		fmt.Printf("first room creation in progress..\n")
		chat.Rooms = make(map[string]*ChatRoom, 5)
	}
	if _, exists := chat.Rooms[roomName]; exists {
		return errors.New("Name already in use\n")
	}
	room := NewChatRoom()
	chat.Rooms[roomName] = &room
	fmt.Printf("Rooms are: %v\n", chat.GetRoomsList())
	return nil
}

func (chat Chat) GetRoomsList() []string {
	rooms := make([]string, len(chat.Rooms))
	i := 0
	for r := range chat.Rooms {
		rooms[i] = r
		i++
	}
	return rooms
}

func (chat *Chat) JoinRoom(roomName string, user *user.User) error {
	room, exists := chat.Rooms[roomName]
	if !exists {
		rooms := chat.GetRoomsList()
		msg := "Room " + roomName + " does not exist\n"
		msg += "Rooms: " + strings.Join(rooms, ",") + "("
		msg += strconv.Itoa(len(chat.Rooms)) + ")\n"
		return errors.New(msg)
	}
	room.Join(user)
	return nil
}

func (chat *Chat) SendMessage(user *user.User, roomNameAndMessage string) error {
	roomNameAndMessageS := string(roomNameAndMessage)
	for roomName := range chat.Rooms {
		if strings.HasPrefix(roomNameAndMessageS, roomName) {
			msg := string(roomNameAndMessage[len(roomName):][:])
			fmt.Printf("%v %v\n", len(msg), msg)
			chat.sendMessage(user, roomName, []byte(msg))
		}
	}
	return nil
}

func (chat *Chat) sendMessage(user *user.User, roomName string, message []byte) error {
	room, exists := chat.Rooms[roomName]
	if !exists {
		return errors.New("Room " + roomName + " does not exist")
	}
	i := 0
	msg := string(message[:])
	fmt.Printf("Sending to room: %v\n", room)
	for i < len(room.Chatters) {
		chatter := room.Chatters[i]
		i++
		if chatter == nil {
			continue
		}
		if chatter.Name == user.Name {
			fmt.Printf("Skipping: %v\n", chatter.Name)
			continue
		}
		fmt.Printf("Sending: %v\n", msg)
		chatter.Connection.Write([]byte(user.Name + ": " + msg))
	}
	return nil
}

func (chat *Chat) HasUser(userName string) bool {
	for room := range chat.Rooms {
		users := chat.Rooms[room].Chatters
		for j := 0; j < len(users); j++ {
			if userName == users[j].Name {
				return true
			}
		}
	}
	return false
}

func (chat *Chat) RemoveUser(connection net.Conn) {
	fmt.Printf("Searching for user to remove...\n")
	for room := range chat.Rooms {
		users := chat.Rooms[room].Chatters
		fmt.Printf("Users (room): %v (%v)\n", users, room)
		for j := 0; j < len(users); j++ {
			if connection == users[j].Connection {
				fmt.Printf("Found user to remove\n")
				newUsers := make([]*user.User, 0)
				if j+1 == len(users) {
					newUsers = append(newUsers, users[:j]...)
				} else {
					newUsers = append(users[:j], users[j+1:]...)
				}
				chat.Rooms[room].Chatters = newUsers
				fmt.Printf("after removal: %v\n", newUsers)
				fmt.Printf("after removal: %v\n", chat.Rooms[room].Chatters)
			}
		}
	}
}
