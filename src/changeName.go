package backend

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/mattn/go-sqlite3"
)

// ChangeDisplayName retrieves the user ID from the SESSION table (by sessionID)
// and updates that user's name in the USER table.
func ChangeDisplayName(sessionID int, newName string) error {
    //  Open the SQLite database.
    db, err := sql.Open("sqlite3", "../sql/BrickedUpDatabase.sql")
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }
    defer db.Close()

    // Look up the userID in the SESSION table.
    var userID int
    err = db.QueryRow("SELECT userid FROM SESSION WHERE id = ?", sessionID).Scan(&userID)
    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("no session found for session ID %d", sessionID)
        }
        return fmt.Errorf("failed to retrieve user ID from session: %w", err)
    }

    // Update the user’s display name in the USER table.
    query := "UPDATE USER SET name = ? WHERE id = ?"
    result, err := db.Exec(query, newName, userID)
    if err != nil {
        return fmt.Errorf("failed to update display name: %w", err)
    }

    // Check if the update actually affected any row.
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get affected rows: %w", err)
    }
    if rowsAffected == 0 {
        return fmt.Errorf("no rows updated; user with ID %d may not exist", userID)
    }

    // Log a success message and return.
    log.Printf("Successfully updated display name for user %d to '%s'\n", userID, newName)
    return nil
}
