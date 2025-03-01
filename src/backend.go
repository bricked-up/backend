// Package backend provides the backend infrastructure (route handling + database)
// for the Bricked-Up website.
package backend

import (
	"fmt"
	"net/http"
)

// Endpoints maps URL paths to their corresponding handler functions.
var endpoints = map[string]http.HandlerFunc{
	"/login":  loginHandler,
	"/signup": signupHandler,
	"/verify": verifyHandler,
}

// LoginHandler handles requests to the /login endpoint.
// It only allows GET requests and responds with a placeholder message.
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method unsupported", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "TODO: Login")
}

// SignupHandler handles requests to the /signup endpoint.
// It restricts the request method to GET and responds with a placeholder message.
func signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method unsupported", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "TODO: Signup")
}

// VerifyHandler handles requests to the /verify endpoint.
// Only GET requests are supported, and it returns a placeholder response.
func verifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method unsupported", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "TODO: Verify")
}

// MainHandler checks if the request URL matches a known endpoint.
// If it does, the corresponding handler is called; otherwise, it returns a 404 error.
func MainHandler(w http.ResponseWriter, r *http.Request) {
	if handler, ok := endpoints[r.URL.Path]; ok {
		handler(w, r)
		return
	}
	http.NotFound(w, r)
}
