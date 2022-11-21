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
	Attempt       int
}

var InputChan = make(chan string, 1)
var ResponseChan = make(chan string, 1)
var LevelChan = make(chan string, 1)
var AttemptChan = make(chan int, 1)

var Data Hangman

func hangHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/hangmanweb.html")
	if r.Method == "GET" {
		level := r.FormValue("lvl")
		fmt.Println(level)
		LevelChan <- level
	}
	fmt.Println("Je vais parser la page")
	if r.Method == "POST" {
		InputChan <- r.FormValue("input")
		Data.WordToDisplay = <-ResponseChan
		Data.Attempt = <-AttemptChan
	}
	fmt.Println(Data.WordToDisplay)
	Data.WordToDisplay = <-ResponseChan
	InputChan <- r.FormValue("input")
	Data.Attempt = <-AttemptChan
	fmt.Println("Je vais afficher la page")
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
	go hangman.Hangman(InputChan, ResponseChan, LevelChan, AttemptChan)
	wg.Wait()
}
