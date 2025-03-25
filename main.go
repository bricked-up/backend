package main

import (
	"brickedup/backend"
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"
)

const PORT = ":3100" // Server port number.
// Main sets up the server and starts listening on the defined PORT.
// It registers MainHandler to process all incoming HTTP requests.
func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

        db, err := sql.Open("sqlite", os.Getenv("DB"))
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            log.Panic(err)
        }

        defer db.Close()
        backend.MainHandler(db, w, r);
    })


	log.Printf("Listening on localhost%s", PORT)
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatal(err)
	}
}
