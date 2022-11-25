package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"text/template"

	"github.com/TomJegou/hangman-classic-Remy/src"
)

// Structures
type User struct {
	Name   string
	Passwd string
	Points int
}

type UserList struct {
	List map[string]User
}
type Hangman_Data struct {
	WordToDisplay, Level, Word, ErrorLogin string
	Points, TotalPoints, Attempt           int
	UsedLetters                            []string
}

// Channels
var InputChan = make(chan string, 1)
var ResponseChan = make(chan string, 1)
var LevelChan = make(chan string, 1)
var AttemptChan = make(chan int, 1)
var WordChan = make(chan string, 1)
var QuitChan = make(chan bool, 1)
var UsedLettersChan = make(chan []string, 1)

// Global Variables
var Current_User User
var Data Hangman_Data
var User_list UserList
var levelHandlerRequestCount int = 0

// Functions Handlers
func hangHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Je vais parser la page")
	t, _ := template.ParseFiles("static/html/hangmanweb.html")
	if r.Method == "GET" {
		Data.Level = r.FormValue("lvl")
		LevelChan <- Data.Level
		Data.WordToDisplay = <-ResponseChan
		InputChan <- "b0c9713aa009f4fcf39920d0d7eda80714b0c44ff2f98205278be112c755ca45e5386cbe7a9fca360ad22f06e45f80a8b8f23838725d15f889e202f5cea26359"
		Data.UsedLetters = <-UsedLettersChan
		InputChan <- "b0c9713aa009f4fcf39920d0d7eda80714b0c44ff2f98205278be112c755ca45e5386cbe7a9fca360ad22f06e45f80a8b8f23838725d15f889e202f5cea26359"
		Data.Attempt = <-AttemptChan
		t.Execute(w, Data)
	} else if r.Method == "POST" {
		InputChan <- r.FormValue("input")
		Data.UsedLetters = <-UsedLettersChan
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
		username := r.FormValue("username")
		password := r.FormValue("password")
		loadUserList()
		if UserExists(User_list.List, username) {
			if User_list.List[username].Passwd == password {
				Current_User.Name = username
				Current_User.Passwd = password
				Current_User.Points = User_list.List[username].Points
				levelHandlerRequestCount = 1
			} else {
				Data.ErrorLogin = "Wrong Password"
				levelHandlerRequestCount = 0
				http.Redirect(w, r, "/login", http.StatusFound)
			}
		} else {
			Data.ErrorLogin = "User don't exists"
			levelHandlerRequestCount = 0
			http.Redirect(w, r, "/login", http.StatusFound)
		}
		t, _ := template.ParseFiles("static/html/ChoiceLvl.html")
		t.Execute(w, Data)
		if levelHandlerRequestCount == 1 {
			go src.Hangman(InputChan, ResponseChan, LevelChan, AttemptChan, WordChan, QuitChan, UsedLettersChan)
		} else {
			QuitChan <- true
			go src.Hangman(InputChan, ResponseChan, LevelChan, AttemptChan, WordChan, QuitChan, UsedLettersChan)
		}
		levelHandlerRequestCount++
	}
}

func winHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/html/win.html")
	Data.Points = src.Points(Data.Attempt, Data.Level)
	Data.TotalPoints += Data.Points
	t.Execute(w, Data)
}

func loseHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/html/lose.html")
	t.Execute(w, Data)
}

func menuHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/html/menu.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, Data)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/html/login.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, Data)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	request := r.FormValue("register")
	if request == "continuez en tant qu'invité" {
		Data.TotalPoints = 0
		http.Redirect(w, r, "/level", http.StatusFound)
	} else {
		fmt.Println(request)
		t, err := template.ParseFiles("static/html/register.html")
		if err != nil {
			log.Fatal(err)
		}
		t.Execute(w, Data)
	}
}

func registerOperationHandler(w http.ResponseWriter, r *http.Request) {
	Current_User.Name = r.FormValue("username")
	Current_User.Passwd = r.FormValue("password")
	Current_User.Points = 0
	User_list.List[Current_User.Name] = Current_User
	saveUserList()
	http.Redirect(w, r, "/login", http.StatusFound)
}

// Functions
func StartServer(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("The server is Running")
	fmt.Println("http://localhost:8080/menu")
	fs := http.FileServer(http.Dir("./static"))
	http.HandleFunc("/registeroperation", registerOperationHandler)
	http.HandleFunc("/hangman", hangHandler)
	http.HandleFunc("/level", levelHandler)
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

func saveUserList() {
	bytevalue, err := json.MarshalIndent(User_list, "", "	")
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile("db/accounts.json", bytevalue, 0644)
}

func loadUserList() {
	userListTmp := &UserList{}
	bytevalue, err := os.ReadFile("db/accounts.json")
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(bytevalue, &userListTmp)
	for key, element := range userListTmp.List {
		User_list.List[key] = element
	}
}

func UserExists(userlist map[string]User, username string) bool {
	for key := range userlist {
		if key == username {
			return true
		}
	}
	return false
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go StartServer(&wg)
	wg.Wait()
}
