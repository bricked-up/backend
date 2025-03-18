package backend

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	_ "modernc.org/sqlite"
)

// DeleteTag removes a tag from the database if the user (linked by sessionID) has write permissions.
func DeleteTag(db *sql.DB, sessionID int, tagID int) error {
	var canWrite bool
	var userID int

	// First, check if the provided sessionID exists in the SESSION table.
	// If no matching row is found, the session is considered invalid.
	err := db.QueryRow("SELECT userid FROM SESSION WHERE id = ?", sessionID).Scan(&userID)
	if err != nil {
		// If no session row is found, return a descriptive error.
		if err == sql.ErrNoRows {
			return errors.New("no session found for session ID " + strconv.Itoa(sessionID))
		}
		// For any other database error, return it as-is.
		return err
	}

	// Next, check if the user (linked to the session) has the can_write permission for the specified tag.
	// This query joins multiple tables to ensure the user is part of the project and has a role allowing writes.
	query := `
        SELECT pr.can_write
        FROM TAG t
        JOIN PROJECT_MEMBER pm
            ON t.projectid = pm.projectid
        JOIN PROJECT_MEMBER_ROLE pmr
            ON pm.id = pmr.memberid
        JOIN PROJECT_ROLE pr
            ON pmr.roleid = pr.id
        JOIN SESSION s
            ON pm.userid = s.userid
        WHERE t.id = ?
          AND s.id = ?;
    `
	err = db.QueryRow(query, tagID, sessionID).Scan(&canWrite)
	if err != nil {
		// If there is no row matching this tag and session, it implies the user has no access or the session is invalid for that tag.
		if err == sql.ErrNoRows {
			return errors.New("no matching row => session or permission not found")
		}
		// Return any other database-related error.
		return err
	}

	// If canWrite is false, the user does not have write privileges.
	if !canWrite {
		return fmt.Errorf("user does not have write permissions")
	}

	// If canWrite is true, remove the specified tag from the TAG table.
	deleteQuery := `DELETE FROM TAG WHERE id = ?;`
	_, err = db.Exec(deleteQuery, tagID)
	if err != nil {
		return err
	}

	// On successful deletion, log a confirmation message.
	fmt.Println("Tag deleted")
	return nil
}
