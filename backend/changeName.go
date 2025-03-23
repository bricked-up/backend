package backend

import (
	"brickedup/backend/utils"
	"database/sql"
	"errors"
	"log"
	"strconv"

	_ "modernc.org/sqlite"
)

// ChangeDisplayName retrieves the user ID from the SESSION table (by sessionID)
// and updates that user's name in the USER table without printing errors to stdout.
func ChangeDisplayName(db *sql.DB, sessionID int, newName string) error {
    // Sanitize newName 
    sanitized := utils.SanitizeText(newName, utils.TEXT)
    
    // Look up the userID in the SESSION table.
    var userID int
    err := db.QueryRow("SELECT userid FROM SESSION WHERE id = ?", sessionID).Scan(&userID)
    if err != nil {
        // If no row is found, return a custom error message.
        if err == sql.ErrNoRows {
            return errors.New("no session found for session ID " + strconv.Itoa(sessionID))
        }
        // Otherwise, return the original error from the DB.
        return err
    }

    // Update the user’s display name in the USER table.
    query := "UPDATE USER SET name = ? WHERE id = ?"
    result, err := db.Exec(query, sanitized, userID)
    if err != nil {
        return err
    }

    // Check if the update actually affected any row.
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if rowsAffected == 0 {
        return errors.New("no rows updated; user with ID " + strconv.Itoa(userID) + " may not exist")
    }

    // Log a success message (not an error).
    log.Printf("Successfully updated display name for user %d to '%s'\n", userID, sanitized)
    return nil
}
