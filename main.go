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

	// Persistent log file
	logFilePath := os.Getenv("LOGS")
	if logFilePath != "" {
		logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Panic(err)
		}

		defer logFile.Close()

		log.SetOutput(logFile)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		db, err := sql.Open("sqlite", os.Getenv("DB"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Panic(err)
		}

		defer db.Close()
		backend.MainHandler(db, w, r)
	})

	log.Printf("Listening on localhost%s", PORT)
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatal(err)
	}
}
