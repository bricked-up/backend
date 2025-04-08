package users

import (
	"brickedup/backend/utils"
	"database/sql"
	"errors"
	"strconv"

	_ "modernc.org/sqlite"
)

// UpdateUser retrieves the user ID from the SESSION table (by sessionID)
// and updates the user based on the new values provided.
func UpdateUser(db *sql.DB, sessionID int, user *utils.User) error {
	// Sanitize newName
	user.Name 		= utils.SanitizeText(user.Name, utils.TEXT)
	user.Email 		= utils.SanitizeText(user.Email, utils.EMAIL)
	user.Password 	= utils.SanitizeText(user.Password, utils.PASSWORD)

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
	query := `
	UPDATE USER 
	SET name = ?, avatar = ?, email = ?, password = ?
	WHERE id = ?
	`
	_, err = db.Exec(
		query, 
		user.Name,
		user.Avatar,
		user.Email,
		user.Password,
		userID)

	if err != nil {
		return err
	}

	return nil
}
