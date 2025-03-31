package backend

import (
	"brickedup/backend/utils"
	"database/sql"
	"errors"
	"os"
)

// UpdateUserAvatar updates the avatar path for a verified user based on their session.
func UpdateUserAvatar(db *sql.DB, sessionID int, avatarFilePath string) error {
	// Sanitize avatar path input
	avatarFilePath = utils.SanitizeText(avatarFilePath, utils.TEXT)

	// Check if file exists on disk
	if _, err := os.Stat(avatarFilePath); os.IsNotExist(err) {
		return errors.New("avatar file does not exist")
	}

	// Get user ID from session
	var userID int
	err := db.QueryRow(`SELECT userid FROM SESSION WHERE id = ?`, sessionID).Scan(&userID)
	if err != nil {
		return errors.New("invalid session ID")
	}

	// Check if user is verified
	var verified bool
	err = db.QueryRow(`SELECT verified FROM USER WHERE id = ?`, userID).Scan(&verified)
	if err != nil || !verified {
		return errors.New("user is not verified")
	}

	// Ensure avatar path is unique
	var exists bool
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM USER WHERE avatar = ?)`, avatarFilePath).Scan(&exists)
	if err != nil {
		return errors.New("failed to check avatar uniqueness")
	}
	if exists {
		return errors.New("avatar path is already in use")
	}

	// Update avatar path in DB
	_, err = db.Exec(`UPDATE USER SET avatar = ? WHERE id = ?`, avatarFilePath, userID)
	if err != nil {
		return errors.New("failed to update avatar")
	}

	return nil
}
