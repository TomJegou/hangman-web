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
	fmt.Println(r.Method)
	fmt.Println("Je vais parser la page")
	t, _ := template.ParseFiles("static/hangmanweb.html")
	if r.Method == "GET" {
		level := r.FormValue("lvl")
		fmt.Println(level)
		LevelChan <- level
		Data.WordToDisplay = <-ResponseChan
		InputChan <- "b0c9713aa009f4fcf39920d0d7eda80714b0c44ff2f98205278be112c755ca45e5386cbe7a9fca360ad22f06e45f80a8b8f23838725d15f889e202f5cea26359"
		Data.Attempt = <-AttemptChan
		t.Execute(w, Data)
	} else if r.Method == "POST" {
		InputChan <- r.FormValue("input")
		Data.WordToDisplay = <-ResponseChan
		fmt.Println(Data.WordToDisplay)
		Data.Attempt = <-AttemptChan
		if Data.WordToDisplay == "662fae5f621abdad32655f00103d88d3fc45f2bb" {
			fmt.Println("Win")
			InputChan <- "Win"
			//http.Redirect(w, r, "/", http.StatusFound)
		} else if Data.WordToDisplay == "8df6be46fc07d973c70580c412430566b4d624a8" {
			fmt.Println("Lose")
			InputChan <- "Lose"
			//http.Redirect(w, r, "/", http.StatusFound)
		}
		t.Execute(w, Data)
	}
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
