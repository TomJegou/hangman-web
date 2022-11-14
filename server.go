package main

import (
	"html/template"
	"log"
	"net/http"
)

type Hangman struct {
	WordToDisplay string
}

func Handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("static/début.html"))
	data := Hangman{
		WordToDisplay: "hangman",
	}
	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", Handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
