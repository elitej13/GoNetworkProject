package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	termbox "github.com/nsf/termbox-go"
)

type position struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

var gameClients = make(map[*websocket.Conn]bool)
var gameBroadcast = make(chan position)
var gameUpgrader = websocket.Upgrader{}

func handlePosition() {
	for {
		//Gets the position from the other thread
		pos := <-gameBroadcast
		log.Printf("Sending new position (%d, %d)", pos.X, pos.Y)
		//Iterates through the active clients to send position
		for client := range gameClients {
			//JSONifies the message for passing to javascript
			bytes, err1 := json.Marshal(pos)
			msg := string(bytes)
			//Sends json to client
			err2 := client.WriteJSON(msg)
			if err1 != nil || err2 != nil {
				log.Printf("Error writing json in game:\n%v\n%v", err1, err2)
				client.Close()
				delete(gameClients, client)
			}
		}
	}
}
func listenForGameMessages(ws *websocket.Conn) {
	for {
		var msg position
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading json:\n%v", err)
			delete(gameClients, ws)
			break
		}
		gameBroadcast <- msg
	}
}

// Handles a new game client connection
func handleGameConnections(w http.ResponseWriter, r *http.Request) {
	log.Printf("Incoming game client.")
	//Upgrades the socket to be a websocket
	ws, err := gameUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	//Sets the socket to be closed upon termination
	defer ws.Close()
	//Stores the socket as being connected
	gameClients[ws] = true
	listenForGameMessages(ws)
}

func handleInput() {
	log.Printf("Starting input listener...")
	pos := position{X: 100, Y: 100}
	go handlePosition()
	for {
		//Listens for any of the keyboard events that are captured by termbox
		select {
		case ev := <-events:
			if ev.Type == termbox.EventKey {
				if ev.Key == termbox.KeyArrowUp {
					pos.Y -= 10
					gameBroadcast <- pos
				}
				if ev.Key == termbox.KeyArrowDown {
					pos.Y += 10
					gameBroadcast <- pos
				}
				if ev.Key == termbox.KeyArrowRight {
					pos.X += 10
					gameBroadcast <- pos
				}
				if ev.Key == termbox.KeyArrowLeft {
					pos.X -= 10
					gameBroadcast <- pos
				}
			}

		}
	}
}
