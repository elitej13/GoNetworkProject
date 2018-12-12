package main

// Input - https://github.com/nsf/termbox-go
// Websockets - https://github.com/gorilla/websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	termbox "github.com/nsf/termbox-go"
)

var events = make(chan termbox.Event, 1000)
var upgrader = websocket.Upgrader{}

func main() {
	//Starts a fileserver that serves the public directory as root on the website
	fs := http.FileServer(http.Dir("../Public"))
	http.Handle("/", fs)

	//Starts the two applications
	startChat()
	startGame()

	//Opens the websocket and begins serving
	log.Println("Starting server at localhost on port 9090")
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func startChat() {
	//Starts two function handlers for use of communication to the javascript websockets
	http.HandleFunc("/chat-ws", handleChatConnections)
	//Starts the message handler for the chat clients
	go handleMessages()
}

func startGame() {
	//Starts the game function handler for use of communication to the javascript websockets
	http.HandleFunc("/game-ws", handleGameConnections)
	//Initializes the input
	err := termbox.Init()
	if err != nil {
		//Unable to capture input, game will not work
		log.Fatal(err)
	}
	//Starts the input polling
	go pollEvents()
	//Starts the input handling
	go handleInput()
}
func pollEvents() {
	//Polls events from termbox-go and pushes them into the events channel
	for {
		events <- termbox.PollEvent()
	}
}
