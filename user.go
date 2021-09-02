package main

import (
	"net/http"
)

type User struct {
	Email        string   `json:"email" pg:",pk"`
	Password     string   `json:"password,omitempty"`
	FirstName    string   `json:"firstName"`
	LastName     string   `json:"lastName"`
	Balance      int64    `json:"balance"`
	Subscription []string `json:"subscription"`
	Shares       []Share  `json:"shares"`
}

type Share struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

func LoginUser(w http.ResponseWriter, r *http.Request) {

}

func SignupUser(w http.ResponseWriter, r *http.Request) {

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
