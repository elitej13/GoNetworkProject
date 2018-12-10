package main

import (
	"log"
	termbox "github.com/nsf/termbox-go"
	"github.com/gorilla/websocket"
)

type position struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

var positions = make(map[*websocket.Conn]position)

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

func handleInput() {
	var pos = position(X: 100, Y: 100)
	for {
		// if any of the keyboard events are captured
		select {
		case ev := <-events:
			if ev.Type == termbox.EventKey {

				if ev.Key == termbox.KeyArrowUp {
					pos.Y += 10
				}
				if ev.Key == termbox.KeyArrowDown {
					pos.Y -= 10
				}
				if ev.Key == termbox.KeyArrowRight {
					pos.X += 10
				}
				if ev.Key == termbox.KeyArrowLeft {
					pos.X -= 10
				}
			}


		}
	}
}
