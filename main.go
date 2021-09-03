package main

import (
	// "fmt"
	"github.com/go-pg/pg/v10"
	// "github.com/go-pg/pg/v10/orm"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var db *pg.DB

func HandleRequests() {
	router := mux.NewRouter()
	router.HandleFunc("/signup", SignupUser).Methods("POST")
	router.HandleFunc("/login", LoginUser).Methods("POST")
	router.HandleFunc("/logout", LogoutUser).Methods("POST")
	router.HandleFunc("/user/buyshare", BuyShare).Methods("POST")
	router.HandleFunc("/user/sellshare", SellShare).Methods("POST")
	router.HandleFunc("/user/subscribe", Subscribe).Methods("POST")
	router.HandleFunc("/user/unsubscribe", Unsubscribe).Methods("POST")

	router.HandleFunc("/user/addbalance", AddBalance).Methods("PUT")

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

	//creates table users, ran once to initialize table
	// err := db.Model((*User)(nil)).CreateTable(&orm.CreateTableOptions{})
	// if err != nil {
	// 	fmt.Println("Table exists", err)
	// }

	HandleRequests()
}
