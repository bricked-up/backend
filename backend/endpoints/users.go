package endpoints

import (
	"brickedup/backend/users"
	"brickedup/backend/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// LoginHandler handles POST requests to the user logins on /login.
func LoginHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method unsupported", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")

	sessionid, err := users.Login(db, email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	session := fmt.Sprint(sessionid)

	cookie := &http.Cookie{
		Name: 		LoginCookie,
		Value:    	session,
		Expires:  	time.Now().Add(12 * 30 * 24 * time.Hour),
		Secure:   	true,
		HttpOnly: 	true,
	}

	http.SetCookie(w, cookie)

}

// SignupHandler handles POST requests for user singups on /signup.
func SignupHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method unsupported", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")

	err := users.Signup(db, email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

}

// VerifyHandler handles GET requests to verify the user email on /verify.
// Takes in `code` as a URL parameter.
func VerifyHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method unsupported", http.StatusMethodNotAllowed)
		return
	}

	code := r.URL.Query().Get("code")

	if code == "" {
		http.Error(w, "Missing id or code", http.StatusBadRequest)
		return
	}

	err := users.VerifyUser(code, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	dashboard := fmt.Sprintf("http://%s:3000/dashboard", os.Getenv("HOST"))
	http.Redirect(w, r, dashboard, http.StatusFound)
}

// UpdateUserHandler handles PATCH requests to update the 
// logged-in user's information on /update-user.
func UpdateUserHandler (db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w,"Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()

	if err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie(LoginCookie)
	if err != nil {
		http.Error(w, "Invalid cookie session", http.StatusUnauthorized)
		return
	}

	sessionID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		http.Error(w,"Invalid session ID", http.StatusBadRequest)
		return
	}

	var user utils.User

	user.Name = r.FormValue("name")
	user.Email = r.FormValue("email")
	user.Password = r.FormValue("password")
	user.Avatar = r.FormValue("avatar")

	err = users.UpdateUser(db, sessionID, &user) 
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetUserHandler handles GET requests to retrieve information 
// about a user on /get-user.
// It takes `userid` as a URL parameter.
func GetUserHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
        http.Error(w, "Invalid form data", http.StatusBadRequest)
        return
    }

	userParam := r.URL.Query().Get("userid")
	userid, err := strconv.Atoi(userParam)

	if err != nil {
        http.Error(w, "Invalid parameter for userid", http.StatusBadRequest)
        return
	}


	user, err := users.GetUser(db, userid)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(user)
}

// GetAllUsersHandler handles GET requests to retrieve all verified users. 
// on /get-all-users.
func GetAllUsersHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	userids, err := users.GetAllUsers(db)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("GetAllUsers(): %s\n", err.Error())
		return
	}

	json, err := json.Marshal(userids)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("json.Marshal(): %s\n", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// DeleteUserHandler handles DELETE requests to 
// delete the currently logged-in user on /delete-user.
// It reads the session cookie, validates it, and calls DeleteUser to remove all user data.
func DeleteUserHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
        http.Error(w, "Invalid form data", http.StatusBadRequest)
        return
    }

	// Get session cookie
	cookie, err := r.Cookie(LoginCookie)
	if err != nil {
		http.Error(w, "Missing session cookie", http.StatusUnauthorized)
		return
	}

	sessionID := cookie.Value

	// Call backend logic to delete the user
	if err := users.DeleteUser(db, sessionID); err != nil {
		http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
		log.Println("deleteUser error:", err)
		return
	}

	// Invalidate the cookie
	cleared := &http.Cookie{
		Name:     LoginCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // deletes the cookie
		HttpOnly: true,
	}
	http.SetCookie(w, cleared)

	// Respond with success
	w.WriteHeader(http.StatusOK)
}

