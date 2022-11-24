package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"text/template"

	"github.com/TomJegou/hangman-classic-Remy/src"
)

type Hangman_Data struct {
	WordToDisplay string
	Attempt       int
	Points        int
	Level         string
	Word          string
}

var InputChan = make(chan string, 1)
var ResponseChan = make(chan string, 1)
var LevelChan = make(chan string, 1)
var AttemptChan = make(chan int, 1)
var WordChan = make(chan string, 1)
var QuitChan = make(chan bool, 1)

var Data Hangman_Data

var levelHandlerRequestCount int = 1

func hangHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Je vais parser la page")
	t, _ := template.ParseFiles("static/html/hangmanweb.html")
	if r.Method == "GET" {
		Data.Level = r.FormValue("lvl")
		LevelChan <- Data.Level
		LevelChan <- Data.Level
		Data.WordToDisplay = <-ResponseChan
		InputChan <- "b0c9713aa009f4fcf39920d0d7eda80714b0c44ff2f98205278be112c755ca45e5386cbe7a9fca360ad22f06e45f80a8b8f23838725d15f889e202f5cea26359"
		Data.Attempt = <-AttemptChan
		t.Execute(w, Data)
	} else if r.Method == "POST" {
		InputChan <- r.FormValue("input")
		Data.WordToDisplay = <-ResponseChan
		Data.Attempt = <-AttemptChan
		if Data.WordToDisplay == "50536101b1c465eafbecc8fca26eeb18a2ac8a2f83570bade315c5a112363cdfd820acad2ab234f91d43f0db8fed0cec400a1109ad8f99c21b5b74f59e8bb00d" {
			fmt.Println("Win")
			http.Redirect(w, r, "/win", http.StatusFound)
		} else if Data.WordToDisplay == "889ce65f137b3b9aa1005f417d7972c948b8bb6360cbdd4118cb05a29d37905744fc0dbc3d17c1de02689d837bfea5bb8114a994f9c1a53dddb993139ab2974c" {
			fmt.Println("Lose")
			Data.Word = <-WordChan
			http.Redirect(w, r, "/lose", http.StatusFound)
		}
		t.Execute(w, Data)
	}
}

func levelHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("static/html/ChoiceLvl.html")
		t.Execute(w, Data)
		if levelHandlerRequestCount == 1 {
			go src.Hangman(InputChan, ResponseChan, LevelChan, AttemptChan, WordChan, QuitChan)
		} else {
			QuitChan <- true
			go src.Hangman(InputChan, ResponseChan, LevelChan, AttemptChan, WordChan, QuitChan)
		}
		levelHandlerRequestCount++
	}
}

func winHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/html/win.html")
	Data.Points = src.Points(Data.Attempt, Data.Level)
	t.Execute(w, Data)
}

func loseHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/html/lose.html")
	t.Execute(w, Data)
}

func StartServer(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("The server is Running")
	fmt.Println("http://localhost:8080/level")
	fs := http.FileServer(http.Dir("./static"))
	http.HandleFunc("/hangman", hangHandler)
	http.HandleFunc("/level", levelHandler)
	http.HandleFunc("/win", winHandler)
	http.HandleFunc("/lose", loseHandler)
	http.Handle("/static/", http.StripPrefix("/static", fs))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go StartServer(&wg)
	wg.Wait()
}
