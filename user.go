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
	//Stores username and password from request body and checks if it is equal with database
	reqBody, _ := ioutil.ReadAll(r.Body)
	var user *User = &User{}
	json.Unmarshal(reqBody, user)

	var userDB *User = &User{}
	json.Unmarshal(reqBody, userDB)

	err := db.Model(userDB).WherePK().Column("email", "password").Select()
	if err != nil {
		panic(err)
	}

	if user.Email == userDB.Email && user.Password == userDB.Password {
		json.NewEncoder(w).Encode(map[string]string{"token": CreateToken(user.Email)})
		fmt.Println("Successfully logged in")
		return

	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect email/password"))
		return
	}
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	c := http.Cookie{
		Name:   "token",
		MaxAge: -1}
	http.SetCookie(w, &c)
	w.Write([]byte("Logged out"))
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
		json.NewEncoder(w).Encode(map[string]string{"token": CreateToken(user.Email)})
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
	//Checks if bearer token exists
	claims := AuthChecker(r.Header["Authorization"], w)
	fmt.Println(claims)

	var user User
	user.Email = fmt.Sprintf("%v", claims["email"])
	err := db.Model(&user).WherePK().Column("email", "first_name", "last_name", "balance", "subscriptions", "shares").Select()
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(map[string]User{"user": user})
}
