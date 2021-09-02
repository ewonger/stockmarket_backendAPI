package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type User struct {
	Email         string         `json:"email" pg:",pk"`
	Password      string         `json:"password,omitempty"`
	FirstName     string         `json:"firstName"`
	LastName      string         `json:"lastName"`
	Balance       int64          `json:"balance"`
	Subscriptions []string       `json:"subscriptions"`
	Shares        map[string]int `json:"shares"`
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	//stores username and password from body
	reqBody, _ := ioutil.ReadAll(r.Body)
	var body map[string]string
	json.Unmarshal(reqBody, &body)

	var user *User = &User{}
	json.Unmarshal(reqBody, user)

	err := db.Model(user).WherePK().Column("email", "password").Select()
	if err != nil {
		panic(err)
	}

	if body["email"] == user.Email && body["password"] == user.Password {
		json.NewEncoder(w).Encode(map[string]string{"token": CreateToken(user.Email)})
		fmt.Println("Successfully logged in")
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect email/password"))
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
	user.Shares = make(map[string]int)
	user.Subscriptions = make([]string, 0)

	_, err := db.Model(user).Insert()
	if err != nil {
		//Email exists
		fmt.Println("Error signing up user. Email already exists")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Email already exists"))

	} else {
		json.NewEncoder(w).Encode(map[string]string{"token": CreateToken(user.Email)})
		fmt.Println("Successfully signed up user")
	}
}

func AddBalance(w http.ResponseWriter, r *http.Request) {
	//Checks if bearer token exists
	claims := AuthChecker(r.Header["Authorization"], w)

	reqBody, _ := ioutil.ReadAll(r.Body)

	//grab balance from email
	var user User
	user.Email = fmt.Sprintf("%v", claims["email"])
	err := db.Model(&user).WherePK().Column("balance").Select()
	if err != nil {
		panic(err)
	}

	//parse amount to be added from body
	var body map[string]int
	json.Unmarshal(reqBody, &body)
	fmt.Println(body)
	if body["addBal"] == 0 {
		fmt.Println("no balance to add")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No balance to add"))
		return
	}

	//update balance
	user.Balance += int64(body["addBal"])
	_, err = db.Model(&user).Column("balance").WherePK().Update()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully updated user")
}

func BuyShare(w http.ResponseWriter, r *http.Request) {
	//Checks if bearer token exists
	claims := AuthChecker(r.Header["Authorization"], w)
	if claims == nil {
		return
	}

	//grabs share list from email
	var user User
	user.Email = fmt.Sprintf("%v", claims["email"])
	err := db.Model(&user).WherePK().Column("shares", "balance").Select()
	if err != nil {
		panic(err)
	}

	//parse added shares from body
	var body map[string]map[string]map[string]int
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &body)

	for name, data := range body["sharesToAdd"] {
		//returns error if not enough balance
		user.Balance -= int64(data["priceCents"] * data["quantity"])
		if user.Balance < 0 {
			fmt.Println("not enough balance")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Not enough balance to buy"))
			return
		}

		if _, ok := user.Shares[name]; ok {
			user.Shares[name] += data["quantity"]
		} else {
			user.Shares[name] = data["quantity"]
		}
	}

	_, err = db.Model(&user).Column("shares", "balance").WherePK().Update()
	if err != nil {
		panic(err)
	}
}

func SellShare(w http.ResponseWriter, r *http.Request) {
	//Checks if bearer token exists
	claims := AuthChecker(r.Header["Authorization"], w)
	if claims == nil {
		return
	}

	//grabs share list from email
	var user User
	user.Email = fmt.Sprintf("%v", claims["email"])
	err := db.Model(&user).WherePK().Column("shares", "balance").Select()
	if err != nil {
		panic(err)
	}

	//parse removed shares from body
	var body map[string]map[string]map[string]int
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &body)

	//remove owned shares and adds balance
	for name, data := range body["sharesToRemove"] {
		user.Balance += int64(data["priceCents"] * data["quantity"])

		if _, ok := user.Shares[name]; ok {
			user.Shares[name] -= data["quantity"]
			if user.Shares[name] == 0 {
				delete(user.Shares, name)
			}
		} else {
			fmt.Println("user does not own listed shares")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("user does not own listed shares"))
			return
		}
	}

	_, err = db.Model(&user).Column("shares", "balance").WherePK().Update()
	if err != nil {
		panic(err)
	}
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
