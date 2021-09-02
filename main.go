package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func HandleRequests() {
	router := mux.NewRouter()
	router.HandleFunc("/signup", SignupUser).Methods("POST")
	router.HandleFunc("/login", LoginUser).Methods("POST")
	router.HandleFunc("/user", AddBalance).Methods("PUT")
	router.HandleFunc("/user", BuyShare).Methods("POST")
	router.HandleFunc("/user", SellShare).Methods("POST")
	router.HandleFunc("/user", Subscribe).Methods("POST")
	router.HandleFunc("/user", getPortfolio).Methods("GET")

	log.Fatal(http.ListenAndServe("localhost:8080", router))
}

func main() {
	HandleRequests()
}
