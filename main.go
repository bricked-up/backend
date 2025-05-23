package main

import (
	"brickedup/backend"
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"
)

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

		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}


		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		db, err := sql.Open("sqlite", os.Getenv("DB"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Panic(err)
		}

		db.SetMaxOpenConns(1)

		defer db.Close()
		backend.MainHandler(db, w, r)
	})

	PORT := os.Getenv("PORT")
	HOST := os.Getenv("HOST")

	log.Printf("Listening on http://%s%s", HOST, PORT)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatal(err)
	}

	// log.Printf("Listening on https://%s%s", HOST, PORT)
	// err := http.ListenAndServeTLS(PORT, "/cert/brickedup.crt", "/cert/priv.key", nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
