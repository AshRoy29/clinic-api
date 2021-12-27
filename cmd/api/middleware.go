package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
)

func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] == nil {
			log.Println("no token found")
			return
		}
		//secretkey := "2dce505d96a53c5768052ee90f3df2055657518dad489160df9913f66042e160"

		tokenString := r.Header["Token"][0]
		log.Println(r.Header["Token"][0])

		mySigningKey := []byte("secret")

		log.Println(mySigningKey)

		log.Println(config.Jwt.Secret)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			log.Println(token.Method.Alg())
			log.Println(token.Method.(*jwt.SigningMethodHMAC))

			log.Println(token.Method)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Println("error in parsing")
				return nil, fmt.Errorf("error in parsing")
			}
			log.Println("OK")
			return mySigningKey, nil
		})

		if err != nil {
			log.Println(err)
			//err := errors.New("your token has expired")
			//json.NewEncoder(w).Encode(err)
			//return
		}

		log.Println(token.Claims.(jwt.MapClaims))
		log.Println(token)
		log.Println(token.Valid)

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			log.Println("vbnm")
			if claims["role"] == "admin" {
				r.Header.Set("Role", "admin")
				handler.ServeHTTP(w, r)
				return
			} else if claims["role"] == "user" {
				r.Header.Set("Role", "user")
				handler.ServeHTTP(w, r)
				return
			}
		}
		errs := errors.New("your token has expired")
		json.NewEncoder(w).Encode(errs)
		json.NewEncoder(w).Encode(err)

	}
}
