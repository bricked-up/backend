package backend

import (
	"brickedup/backend/utils"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"

	_ "modernc.org/sqlite" // SQLite driver for database/sql
)

// GenerateVerificationCode generates a random hex-encoded code
func generateVerificationCode() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// SendVerificationEmail sends an email using gomail with a verification code
func sendVerificationEmail(to string, code string) {
	email := os.Getenv("EMAIL")
    password := os.Getenv("PASS")

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
	sanitizedEmail := utils.SanitizeText(email, utils.EMAIL)
	sanitizedPassword := utils.SanitizeText(password, utils.PASSWORD)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(sanitizedPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Insert user into database
	res, err := db.Exec("INSERT INTO USER (email, password, name) VALUES (?, ?, 'New User')", sanitizedEmail, passwordHash)
	if err != nil {
		return err
	}
	userID, _ := res.LastInsertId()

	// Generate verification code
	code := generateVerificationCode()

	// Insert verification record into VERIFY_USER table
	res, err = db.Exec("INSERT INTO VERIFY_USER (code, expires) VALUES (?, ?)", code, time.Now().Add(24*time.Hour))
	if err != nil {
		return err
	}
	verifyID, _ := res.LastInsertId()

	// Update USER to link verification ID
	_, err = db.Exec("UPDATE USER SET verifyid = ? WHERE id = ?", verifyID, userID)
	if err != nil {
		return err
	}

	// Send verification email
	sendVerificationEmail(email, code)
	return nil
}
