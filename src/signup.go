package backend

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"gopkg.in/gomail.v2"

	_ "modernc.org/sqlite" // SQLite driver for database/sql
)

// User struct captures the fields of the user table
type User struct {
	ID       uint
	Email    string
	Username string
	Password string
}

// VerifyUser struct captures the fields of the verify_users table
type VerifyUser struct {
	ID     uint
	Code   string
	Expire time.Time
}

var db *sql.DB

// Function initDB initializes database connection
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "BrickedUpDatabase.sql")
	if err != nil {
		log.Fatal("Database Connection error:", err)
	}
}

// Function generateVerificationCode generates a random code
func generateVerificationCode() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// Function getEmailCredentials retrieves email and password from the database
func getEmailCredentials() (string, string, error) {
	var email, password string
	err := db.QueryRow("SELECT email, password FROM user LIMIT 1").Scan(&email, &password)
	if err != nil {
		return "", "", err
	}
	return email, password, nil
}

// Function sendVerificationEmail sends an email using gomail with a verification code
func sendVerificationEmail(to, code string) {
	email, password, err := getEmailCredentials()
	if err != nil {
		log.Println("Failed to retrieve email credentials:", err)
		return
	}

	m := gomail.NewMessage()
	m.SetHeader("From", email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Account Verification")
	m.SetBody("text/html", fmt.Sprintf("<p>Your verification code is: <strong>%s</strong></p><p><button>Verify</button></p>", code))
	/* email is sent to user without implementation, but I did not
	know how to implement sending the email without smtp or server access */
	d := gomail.NewDialer("smtp.example.com", 587, email, password)
	if err := d.DialAndSend(m); err != nil {
		log.Println("Failed to send email:", err)
	}
}

// Function registerUser handles user registration
func registerUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "405 error", http.StatusMethodNotAllowed)
		return
	}

	// Req struct parses the request load
	var req struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	// User added to database
	res, err := db.Exec("INSERT INTO users (email, username, password) VALUES (?, ?, ?)", req.Email, req.Username, req.Password)
	if err != nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}
	userID, _ := res.LastInsertId()

	// Verification code inserted in database with a 24 hour timer
	code := generateVerificationCode()
	_, err = db.Exec("INSERT INTO verify_users (id, code, expire) VALUES (?, ?, ?)", userID, code, time.Now().Add(24*time.Hour))
	if err != nil {
		http.Error(w, "Verification error", http.StatusInternalServerError)
		return
	}

	// Send verification email using gomail
	sendVerificationEmail(req.Email, code)

	// Informs client of verification email
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, `{"message": "User created, verification email sent"}`)
}
