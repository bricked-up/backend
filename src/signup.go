package backend

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"

	"fmt"
	"log"
	"time"

	"gopkg.in/gomail.v2"

	_ "modernc.org/sqlite" // SQLite driver for database/sql
)

// GenerateVerificationCode generates a random code
func generateVerificationCode() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// SendVerificationEmail sends an email using gomail with a verification code
func sendVerificationEmail(to, code string) {
	email := "backend@gmail.com"
	password := "123"

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

// RegisterUser handles user registration
func registerUser(db *sql.DB, email, password string) error {

	// User added to database
	res, err := db.Exec("INSERT INTO users (email, password) VALUES (?, ?)", email, password)
	if err != nil {
		return err
	}
	userID, _ := res.LastInsertId()

	// Verification code inserted in database with a 24 hour timer
	code := generateVerificationCode()
	_, err = db.Exec("INSERT INTO verify_users (id, code, expire) VALUES (?, ?, ?)", userID, code, time.Now().Add(24*time.Hour))
	if err != nil {

		return err
	}

	// Send verification email using gomail
	sendVerificationEmail(email, code)
	return nil
}
