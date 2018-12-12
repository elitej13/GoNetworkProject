package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

//Message struct for packaging and passing the
//relevant message data to and from the javascript
type message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
	Color    string `json:"color"`
}

var chatClients = make(map[*websocket.Conn]bool)
var chatBroadcast = make(chan message)

func handleMessages() {
	for {
		//Receives messages from the chat channel and pushes them to other clients
		msg := <-chatBroadcast
		for client := range chatClients {
			//Sends the received message to the client
			err := client.WriteJSON(msg)
			if err != nil {
				//Json communication has goen wrong and will be killed off for the better good
				log.Printf("Error writing json in chat: %v", err)
				client.Close()
				delete(chatClients, client)
			}
		}
	}
}

func listenForChatMessages(ws *websocket.Conn) {
	for {
		//Starts to listen for a message from the given websocket
		var msg message
		err := ws.ReadJSON(&msg)
		log.Printf("[%s] %s", msg.Username, msg.Message)
		if err != nil {
			//Socket has goen wrong and will be killed off for the better good
			log.Printf("Error reading json: %v", err)
			delete(chatClients, ws)
			break
		}
		//Pushes the received message into the chat channel
		chatBroadcast <- msg
	}
}

func handleChatConnections(w http.ResponseWriter, r *http.Request) {
	log.Printf("Incoming chat client!")
	//Upgrades the http socket to a websocket for two way communication
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	//Sets the socket to be closed upon termination
	defer ws.Close()
	//Records the socket as being an active connection
	chatClients[ws] = true
	//Starts listening for messages from the websocket
	listenForChatMessages(ws)
}
