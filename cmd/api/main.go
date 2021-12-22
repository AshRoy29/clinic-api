package main

import (
	"fmt"
	"net/http"
)

func main() {
	r := Routes()
	fmt.Println("Server is getting started...")
	http.ListenAndServe(":4000", r)
	fmt.Println("Listening at port 4000...")
}
