package src

import (
	"encoding/json"
	"log"
	"os"
)

type User struct {
	Name   string
	Passwd string
	Points int
}

type UserList struct {
	List []User
}

var Current_User User
var User_list UserList

func saveUserList() {
	bytevalue, err := json.MarshalIndent(User_list, "", "	")
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile("db/accounts.json", bytevalue, 0644)
}

func loadUserList() {
	bytevalue, err := os.ReadFile("db/accounts.json")
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(bytevalue, &User_list)
}

func UserExists(userlist []User, username string) (bool, int) {
	for index, elements := range userlist {
		if elements.Name == username {
			return true, index
		}
	}
	return false, 0
}

func savePoints() {
	User_list.List[IndexUserList].Points = Current_User.Points
}
