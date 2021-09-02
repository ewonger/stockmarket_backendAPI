package main

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var db *pg.DB

func HandleRequests() {
	router := mux.NewRouter()
	router.HandleFunc("/signup", SignupUser).Methods("POST")
	router.HandleFunc("/login", LoginUser).Methods("POST")
	router.HandleFunc("/user", BuyShare).Methods("POST")
	router.HandleFunc("/user", SellShare).Methods("POST")
	router.HandleFunc("/user", Subscribe).Methods("POST")

	router.HandleFunc("/user", AddBalance).Methods("PUT")

	router.HandleFunc("/user", getPortfolio).Methods("GET")

	log.Fatal(http.ListenAndServe("localhost:8080", router))
}

func main() {
	//connect to db
	db = pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     "postgres",
		Password: "postgres",
		Database: "testdb",
	})
	defer db.Close()

	//creates table users
	err := db.Model((*User)(nil)).CreateTable(&orm.CreateTableOptions{})
	if err != nil {

		fmt.Println("Table exists")
	}

	HandleRequests()
}
