package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func chatPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("chat.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		fmt.Println("Message:", r.Form["Message"])
		fmt.Println("ChatWindow:", r.Form["ChatWindow"])
	}
}

func gamePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("game.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
	}
}

func landingPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("landing.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		fmt.Println(r.Form)
		fmt.Println("path", r.URL.Path)
		fmt.Println("scheme", r.URL.Scheme)
		fmt.Println(r.Form["url_long"])
		for k, v := range r.Form {
			fmt.Println("key:", k)
			fmt.Println("val:", strings.Join(v, ""))
		}
	}
	//fmt.Fprintf(w, "Simple Chat Program") // write data to response
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", landingPage)
	mux.HandleFunc("/chat", chatPage)
	mux.HandleFunc("/game", gamePage)

	fileServer := http.FileServer(http.Dir("./rsc"))
	mux.Handle("/rsc/", http.StripPrefix("/rsc", fileServer))

	log.Println("Starting server on :9090")
	err := http.ListenAndServe(":9090", mux)
	log.Fatal(err)
}
