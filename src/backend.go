// Package backend provides the backend infrastructure (route handling + database)
// for the Bricked-Up website.
// package backend;
package backend

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
)

type DBHandlerFunc func(*sql.DB, http.ResponseWriter, *http.Request)

// Endpoints maps URL paths to their corresponding handler functions.
var endpoints = map[string]DBHandlerFunc{
	"/login":  loginHandler,
	"/signup": signupHandler,
	"/verify": verifyHandler,
}

// LoginHandler handles requests to the /login endpoint.
// It only allows GET requests and responds with a placeholder message.
func loginHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method unsupported", http.StatusMethodNotAllowed)
		return
	}

    r.ParseForm();
    email := r.FormValue("email")
    password := r.FormValue("password")

    sessionid, err := login(db, email, password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        log.Println(err.Error())
        return
    }

    session := fmt.Sprint(sessionid)

    cookie := &http.Cookie{
		Name:    "bricked-up_login",
		Value:   session,
		Expires: time.Now().Add(12 * 30 * 24 * time.Hour),
		Secure:  true,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
}

// SignupHandler handles requests to the /signup endpoint.
// It restricts the request method to GET and responds with a placeholder message.
func signupHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method unsupported", http.StatusMethodNotAllowed)
		return
	}

    r.ParseForm();
    email := r.FormValue("email")
    password := r.FormValue("password")

    err := registerUser(db, email, password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        log.Println(err.Error())
        return
    }
}

// VerifyHandler handles requests to the /verify endpoint.
// Only GET requests are supported, and it returns a placeholder response.
func verifyHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method unsupported", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "TODO: Verify")
}

// MainHandler checks if the request URL matches a known endpoint.
// If it does, the corresponding handler is called; otherwise, it returns a 404 error.
func MainHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if handler, ok := endpoints[r.URL.Path]; ok {
		handler(db, w, r)
		return
	}
	http.NotFound(w, r)
}
