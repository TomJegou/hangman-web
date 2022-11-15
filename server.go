package main

import (
	"fmt"
	"hangman"
	"html/template"
	"log"
	"net/http"
	"sync"
)

type Hangman struct {
	WordToDisplay string
	Method        string
}

var InputToHangman string
var ResponseFromHangman string

func input(wg *sync.WaitGroup, inputChan chan<- string, responseChan <-chan string) {
	fmt.Println("Input routine activated")
	defer wg.Done()
	for content := range responseChan {
		ResponseFromHangman = content
		inputChan <- ResponseFromHangman
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	data := Hangman{
		WordToDisplay: ResponseFromHangman,
		Method:        r.Method,
	}
	t, _ := template.ParseFiles("static/hangmanweb.html")
	t.Execute(w, data)
	if r.Method == "POST" {
		InputToHangman = r.FormValue("input")
	}
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
	inputChan := make(chan string, 1)
	responseChan := make(chan string, 1)
	wg.Add(2)
	go input(&wg, inputChan, responseChan)
	go Server()
	go hangman.Hangman(&wg, inputChan, responseChan)
	wg.Wait()
}
