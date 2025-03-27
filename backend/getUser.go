package backend

import (
	"database/sql"
	"encoding/json"
	"errors"

	_ "modernc.org/sqlite"
)

// User holds the columns returned by the SELECT query (minus the password).
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Verified bool   `json:"verified"`
	Avatar   string `json:"avatar"`
}

// GetUserDetails fetches one user by ID from the DB and returns JSON data.
func getUserDetails(db *sql.DB, userID int) ([]byte, error) {
	// Get exactly one row for the given userID.
	row := db.QueryRow(`SELECT id, name, email, verified, avatar FROM USER WHERE id = ?`, userID)

	var user User
	// Scan fills our user struct with the row's data or returns an error if no row.
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Verified, &user.Avatar); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("UserId not found")
		}
		return nil, err
	}

	// Convert the user struct to JSON.
	jsonUser, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	return jsonUser, nil
}