package main

import (
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

var Input string
var Content string

func input(wg *sync.WaitGroup, inputChan chan<- string, responseChan <-chan string) {
	defer wg.Done()
	for content := range responseChan {
		Content = content
		inputChan <- Input
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	data := Hangman{
		WordToDisplay: Content,
		Method:        r.Method,
	}
	t, _ := template.ParseFiles("static/hangmanweb.html")
	if r.Method == "POST" {
		Input = r.FormValue("input")
		rootHandler(w, r)
	}
	t.Execute(w, data)
}

func Server() {
	fs := http.FileServer(http.Dir("./static"))
	http.HandleFunc("/", rootHandler)
	http.Handle("/static/", http.StripPrefix("/static", fs))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var wg sync.WaitGroup
	inputChan := make(chan string)
	responseChan := make(chan string)
	wg.Add(2)
	go input(&wg, inputChan, responseChan)
	go Server()
	go hangman.Hangman(&wg, inputChan, responseChan)
	wg.Wait()
}
