package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!")
}

func handlerhi(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi Golang")
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/hi", handlerhi)

	port := ":3000"
	fmt.Println("Server is running on port" + port)

	// Start server on port specified above
	log.Fatal(http.ListenAndServe(port, nil))
}
