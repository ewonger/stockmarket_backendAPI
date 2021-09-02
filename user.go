package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type User struct {
	Email         string   `json:"email" pg:",pk"`
	Password      string   `json:"password,omitempty"`
	FirstName     string   `json:"firstName"`
	LastName      string   `json:"lastName"`
	Balance       int64    `json:"balance"`
	Subscriptions []string `json:"subscriptions"`
	Shares        []Share  `json:"shares"`
}

type Share struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

func LoginUser(w http.ResponseWriter, r *http.Request) {

}

func SignupUser(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var user *User = &User{}
	json.Unmarshal(reqBody, user)

	user.Balance = 0
	_, err := db.Model(user).Insert()
	if err != nil {
		//Email exists
		fmt.Println("Error signing up user. Email already exists")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Email already exists"))
		return

	} else {
		json.NewEncoder(w).Encode(map[string]string{"token": "test"})
		fmt.Println("Successfully signed up user")
		return
	}
}

func AddBalance(w http.ResponseWriter, r *http.Request) {

}

func BuyShare(w http.ResponseWriter, r *http.Request) {

}

func SellShare(w http.ResponseWriter, r *http.Request) {

}

func Subscribe(w http.ResponseWriter, r *http.Request) {

}

func getPortfolio(w http.ResponseWriter, r *http.Request) {

}
