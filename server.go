package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Hangman struct {
	WordToDisplay string
}

func hangHandler(w http.ResponseWriter, r *http.Request) {
	data := Hangman{
		WordToDisplay: "hangman"}
	tmpl, err := template.ParseFiles("static/début.html")
	if err != nil {
		fmt.Println(err)
	}
	tmpl.Execute(w, data)
}

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/hangman", hangHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
