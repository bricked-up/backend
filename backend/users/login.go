package users

import (
	"brickedup/backend/utils"
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

// Login authenticates a user by verifying their email and password.
// If authentication is successful and the user is verified, it creates a new session or reuses an existing one.
// It returns the session data.
func Login(db *sql.DB, email, password string) (session *utils.SessionData, err error) {
	session = &utils.SessionData{}
	var storedPassword string

    // Query the database to get the user's ID, hashed password, and verification status
    err = db.QueryRow(
        `SELECT id, password 
		FROM USER 
		WHERE email = ? AND verifyid IS NULL `, 
        email).Scan(&session.UserID, &storedPassword)

	if err != nil {
        return nil, err
    }

    // Compare the provided password with the stored hashed password
    if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)); 
    err != nil {
        return nil, err
    }

    // Set the session expiration time (valid for 24 hours)
    session.Expires = time.Now().Add(24 * time.Hour)

    // Insert the new session into the SESSION table in the database
    result, err := db.Exec(
		"INSERT INTO SESSION (userid, expires) VALUES (?, ?)", 
		session.UserID, session.Expires)

    if err != nil {
        return nil, err
    }

    session.SessionID, err = result.LastInsertId()
    if err != nil {
        return nil, err
    }

    return session, nil
}
