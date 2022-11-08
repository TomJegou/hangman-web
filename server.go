package main

import (
	"log"
	"net/http"
)

func hangHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	//http.HandleFunc("/hangman", )
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
