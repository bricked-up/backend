package users

import (
	"database/sql"
	"fmt"
)

// DeleteUser removes a user and associated records based on the session ID.
func DeleteUser(db *sql.DB, sessionID string) error {
	var userID int

	// Retrieve user ID from session
	err := db.QueryRow("SELECT userid FROM SESSION WHERE id = ?", sessionID).Scan(&userID)
	if err != nil {
		return err
	}

	// Debugging: Log found user ID
	fmt.Printf("Deleting user ID: %d\n", userID)

	// Delete user-related entries in foreign key tables
	_, err = db.Exec("DELETE FROM REMINDER WHERE userid = ?", userID)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM USER_ISSUES WHERE userid = ?", userID)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM PROJECT_MEMBER WHERE userid = ?", userID)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM ORG_MEMBER WHERE userid = ?", userID)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM SESSION WHERE userid = ?", userID)
	if err != nil {
		return err
	}

	// Finally, delete the user
	_, err = db.Exec("DELETE FROM USER WHERE id = ?", userID)
	if err != nil {
		return err
	}

	fmt.Printf("User ID %d successfully deleted\n", userID)
	return nil
}
