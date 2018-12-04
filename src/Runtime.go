package main

// Keyboard input - https://github.com/nsf/termbox-go

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error writing json: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()
	clients[ws] = true

	for {
		var msg Message

		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading json: %v", err)
			delete(clients, ws)
			break
		}
		broadcast <- msg
	}
}

func main() {

	fs := http.FileServer(http.Dir("../Public"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	log.Println("Starting server on :9090")
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal(err)
	}

}
