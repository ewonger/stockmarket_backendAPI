package main

import (
	"crypto/rand"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

var hmacSecret []byte

//Generate random bytes for HMAC
func GenerateSecret(tokenBytes []byte) []byte {
	tokenBytes = make([]byte, 4)
	rand.Read(tokenBytes)
	return tokenBytes
}

//Creates JWT token
func CreateToken(email string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Minute * 5).Unix(),
	})

	tokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		fmt.Println("signed string err")
	}
	return tokenString
}

//Verifies JWT token
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		fmt.Println("Token Valid")
	} else {
		fmt.Println(err)
	}

	return claims, err
}
