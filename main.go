package main

import (
	backend "brickedup/backend/src"
	"log"
	"net/http"
)

const PORT = ":3100" // Server port number.
// Main sets up the server and starts listening on the defined PORT.
// It registers MainHandler to process all incoming HTTP requests.
func main() {
	http.HandleFunc("/", backend.MainHandler)

	log.Printf("Listening on localhost%s", PORT)
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatal(err)
	}
}
