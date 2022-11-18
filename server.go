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
var LevelChan = make(chan string, 1)

var Data Hangman

func hangHandler(w http.ResponseWriter, r *http.Request) {
	LevelChan <- r.FormValue("lvl")
	t, _ := template.ParseFiles("static/hangmanweb.html")
	if r.Method == "POST" {
		InputChan <- r.FormValue("input")
		Data.WordToDisplay = <-ResponseChan
	}
	Data.WordToDisplay = <-ResponseChan
	InputChan <- r.FormValue("input")
	t.Execute(w, Data)
}

func levelHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/startMenu.html")
	t.Execute(w, Data)
}

func Server(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("The server is Running")
	fs := http.FileServer(http.Dir("./static"))
	http.HandleFunc("/hangman", hangHandler)
	http.HandleFunc("/", levelHandler)
	http.Handle("/static/", http.StripPrefix("/static", fs))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go Server(&wg)
	go hangman.Hangman(InputChan, ResponseChan, LevelChan)
	wg.Wait()
}
