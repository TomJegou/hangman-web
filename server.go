package main

import (
	"fmt"
	"hangman"
	"log"
	"net/http"
	"sync"
	"text/template"
)

type Hangman struct {
	WordToDisplay string
}

var InputChan = make(chan string, 1)
var ResponseChan = make(chan string, 1)

var Data Hangman

func rootHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/hangmanweb.html")
	if r.Method == "POST" {
		Data.WordToDisplay = <-ResponseChan
		InputChan <- r.FormValue("input")
	}
	InputChan <- r.FormValue("input")
	Data.WordToDisplay = <-ResponseChan
	t.Execute(w, Data)
}

func Server() {
	fmt.Println("The server is Running")
	fs := http.FileServer(http.Dir("./static"))
	http.HandleFunc("/", rootHandler)
	http.Handle("/static/", http.StripPrefix("/static", fs))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go Server()
	go hangman.Hangman(&wg, InputChan, ResponseChan)
	wg.Wait()
}
