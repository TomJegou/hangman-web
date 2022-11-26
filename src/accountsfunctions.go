package src

import (
	"encoding/json"
	"log"
	"os"
)

// Stuctures
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

/*Save the list of users in the accounts.json file*/
func saveUserList() {
	bytevalue, err := json.MarshalIndent(User_list, "", "	")
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile("db/accounts.json", bytevalue, 0644)
}

/*Load the list of users from the accounts.json file*/
func loadUserList() {
	bytevalue, err := os.ReadFile("db/accounts.json")
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(bytevalue, &User_list)
}

/*
Check if the username exists in the userlist, if yes it returns true and
the index of that user in the list
*/
func UserExists(userlist []User, username string) (bool, int) {
	for index, elements := range userlist {
		if elements.Name == username {
			return true, index
		}
	}
	return false, -1
}

/*Save the Current_User points into his corresponding user data in the User_list*/
func savePoints() {
	User_list.List[IndexUserList].Points = Current_User.Points
}
