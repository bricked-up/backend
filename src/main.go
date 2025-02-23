// Package main provides the backend infrastructure (route handling + database) 
// for the Bricked-Up website.

package main

import (
    "fmt"
    "log"
    "net/http"
)

const PORT = ":3100" // Server port number.

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
func mainHandler(w http.ResponseWriter, r *http.Request) {
    if handler, ok := endpoints[r.URL.Path]; ok {
        handler(w, r)
        return
    }
    http.NotFound(w, r)
}

// Main sets up the server and starts listening on the defined PORT.
// It registers MainHandler to process all incoming HTTP requests.
func main() {
    http.HandleFunc("/", mainHandler)

    log.Printf("Listening on localhost%s", PORT)
    if err := http.ListenAndServe(PORT, nil); err != nil {
        log.Fatal(err)
    }
}
