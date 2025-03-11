package backend

import (
	"database/sql"

	_ "modernc.org/sqlite" // SQLite driver
)

// DeleteUser deletes a user from the users table
func deleteUser(db *sql.DB, sessionid string) error {
	// Retrieve user_id from session table
	var userID string
	err := db.QueryRow("SELECT user_id FROM session WHERE id = ?", sessionid).Scan(&userID)
	if err != nil {
		return err
	}

	// Delete user from users table
	_, err = db.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		return err
	}

	return nil
}
