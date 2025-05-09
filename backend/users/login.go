package users

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

// Login authenticates a user by verifying their email and password.
// If authentication is successful and the user is verified, it creates a new session or reuses an existing one.
// It returns the sessionid.
func Login(db *sql.DB, email, password string) (sessionid int64, err error) {
	var userID int
	var storedPassword string

    // Query the database to get the user's ID, hashed password, and verification status
    err = db.QueryRow(
        `SELECT id, password 
		FROM USER 
		WHERE email = ? AND verifyid IS NULL `, 
        email).Scan(&userID, &storedPassword)

	if err != nil {
        return -1, err
    }

    // Compare the provided password with the stored hashed password
    if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)); 
    err != nil {
        return -1, err
    }

    // Set the session expiration time (valid for 24 hours)
    expiresAt := time.Now().Add(24 * time.Hour)

    // Insert the new session into the SESSION table in the database
    result, err := db.Exec("INSERT INTO SESSION (userid, expires) VALUES (?, ?)", userID, expiresAt)
    if err != nil {
        return -1, err
    }

    sessionid, err = result.LastInsertId()
    if err != nil {
        return -1, err
    }

    return sessionid, nil
}
