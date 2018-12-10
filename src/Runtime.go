package main

// Keyboard input - https://github.com/nsf/termbox-go

import (
	"log"
	"net/http"

	termbox "github.com/nsf/termbox-go"
)

var events = make(chan termbox.Event, 1000)

func main() {

	fs := http.FileServer(http.Dir("../Public"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", handleConnections)

	go handleMessages()
	go pollEvents()
	go handleInput()

	log.Println("Starting server on :9090")
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func pollEvents() {
	for {
		events <- termbox.PollEvent()
	}
}
