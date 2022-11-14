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
	tmpl := template.Must(template.ParseFiles("static/hangmanweb.html"))
	data := Hangman{
		WordToDisplay: Content,
	}
	tmpl.Execute(w, data)
	Input = r.FormValue("input")
}

func Server() {
	http.HandleFunc("/", rootHandler)
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
