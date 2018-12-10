package main

// Keyboard input - https://github.com/nsf/termbox-go

import (
	"log"
	"net/http"
)

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
