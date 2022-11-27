package src

//Imports section
import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

/*Establish all the routing for the web-app and start the server*/
func StartServer(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("The server is Running")
	fmt.Println("http://localhost:8080/menu")
	fs := http.FileServer(http.Dir("./static"))
	http.HandleFunc("/checkCredentials", checkCredentialsHandler)
	http.HandleFunc("/registeroperation", registerOperationHandler)
	http.HandleFunc("/hangman", hangHandler)
	http.HandleFunc("/level", levelHandler)
	http.HandleFunc("/checklevel", checkLevelHandler)
	http.HandleFunc("/win", winHandler)
	http.HandleFunc("/lose", loseHandler)
	http.HandleFunc("/menu", menuHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.Handle("/static/", http.StripPrefix("/static", fs))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
