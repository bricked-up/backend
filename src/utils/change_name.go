package utils

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// ChangeDisplayName updates a user's display name in the USER table based on sessionID.
func ChangeDisplayName(db *sql.DB, sessionID int, newName string) error {
	// Prepare an SQL statement to update the name field for the specified user.
	query := "UPDATE USER SET name = ? WHERE id = ?"

	// Execute the update and handle potential errors.
	result, err := db.Exec(query, newName, sessionID)
	if err != nil {
		return fmt.Errorf("failed to update display name %w", err)
	}

	// Check how many rows were updated to confirm the operation succeeded.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected; session ID might not match")
	}

	// Log a success message when the display name is updated.
	log.Printf("Successfully updated display name for session %d to %s\n", sessionID, newName)
	return nil
}
