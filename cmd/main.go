package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func requestHasBody(r *http.Request) bool {
	if _, hasBody := r.Header["Content-Length"]; !hasBody {
		return false
	}
	return true
}

func getBody(r *http.Request) ([]byte, error) {
	body := r.Body
	requestBody, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return requestBody, nil
}

func landingHandler(w http.ResponseWriter, r *http.Request) {
	if !requestHasBody(r) {
		fmt.Fprintf(w, "No body\n")
		return
	}

	requestBody, err := getBody(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, "No request body")
		fmt.Fprintf(w, string(http.StatusBadRequest))
		return
	}

	fmt.Fprintf(w, string(requestBody))
}

type Chat struct {
	rooms map[string][]string
}

var chat = Chat{}

func joinRoomHandler(w http.ResponseWriter, r *http.Request) {
	if !requestHasBody(r) {
		fmt.Fprintf(w, "No body\n")
		return
	}

	requestBody, err := getBody(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, "No request body")
	}
	roomAndChatter := strings.Split(string(requestBody), "|||")
	if _, exists := chat.rooms[roomAndChatter[0]]; !exists {
		var newRoom = make([]string, 1)
		chat.rooms[roomAndChatter[0]] = newRoom
		fmt.Fprintf(w, "Created room %s\n", roomAndChatter[0])
	}
	fmt.Fprintf(w, "Adding you to room %s\n", roomAndChatter[0])
	chat.rooms[roomAndChatter[0]] = append(chat.rooms[roomAndChatter[0]], roomAndChatter[1])
	fmt.Fprintf(w, "Room Chatters:\n %s\n", strings.Join(chat.rooms[roomAndChatter[0]], "\n"))
	return
}

func main() {
	chat.rooms = make(map[string][]string)
	http.HandleFunc("/", landingHandler)
	http.HandleFunc("/join", joinRoomHandler)
	http.ListenAndServe(":4040", nil)
}
