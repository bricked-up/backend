// Package brickedup/backend provides the backend 
// infrastructure (route handling + database) for the Bricked-Up website.
package main

import (
    "fmt"
    "log"
    "net/http"
)

const Port = ":3100"

//endpoints map
var endpoints = map[string]http.HandlerFunc {
    "/login" : loginHandler,
    "/signup" : signupHandler,
    "/verify" : verifyHandler,
}

//loginHandler, signupHanlder, verifyHandler : helper funcs. for connecting the endpoints 
func loginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method !=  "GET" {
        http.Error(w, "Method unsupported", http.StatusMethodNotAllowed)
        return
    }
    fmt.Fprintf(w, "TODO : login")
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        http.Error(w, "Method unsupported", http.StatusMethodNotAllowed)
        return
    }
    fmt.Fprintf(w, "TODO : signup")
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        http.Error(w, "Method unsupported", http.StatusMethodNotAllowed)
        return
    }
    fmt.Fprintf(w, "TODO : verify")
}

//mainHanlder uses the endpoints map to route requests.
func mainHandler(w http.ResponseWriter, r *http.Request) {
    if handler, ok := endpoints[r.URL.Path] ; ok {
        handler(w,r)
        return
    }
    http.NotFound(w, r)
}

func main() {
    http.HandleFunc ("/", mainHandler)
    
    fmt.Printf("Router starting...")
    if err := http.ListenAndServe(Port, nil); err != nil {
        log.Fatal(err)
    }

}