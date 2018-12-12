package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

var chatClients = make(map[*websocket.Conn]bool)
var chatBroadcast = make(chan message)
var chatUpgrader = websocket.Upgrader{}

func handleMessages() {
	for {
		msg := <-chatBroadcast
		for client := range chatClients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error writing json in chat: %v", err)
				client.Close()
				delete(chatClients, client)
			}
		}
	}
}

func listenForChatMessages(ws *websocket.Conn) {
	for {
		var msg message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading json: %v", err)
			delete(chatClients, ws)
			break
		}
		chatBroadcast <- msg
	}
}

func handleChatConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := chatUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	chatClients[ws] = true
	listenForChatMessages(ws)
}
