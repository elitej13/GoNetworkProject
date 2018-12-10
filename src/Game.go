package main

import (
	"log"
)

type position struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

// var positions = make(map[*websocket.Conn]position)

func handlePosition(pos position) {
	for client := range clients {
		err := client.WriteJSON(pos)
		if err != nil {
			log.Printf("Error writing json in game: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}
