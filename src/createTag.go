package backend

import (
	"brickedup/backend/src/utils"
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

// CreateTag creates a new tag for a project.
// It takes sessionID (int), projectID (int), tagName (string), and tagColor (string) as parameters.
func CreateTag(db *sql.DB, sessionID int, projectID int, tagName string, tagColor string) (int, error) {
	// Validate inputs
	if tagName == "" || tagColor == "" {
		return 0, errors.New("missing tagName or tagColor")
	}

	// Sanitize the tag name and color
	sanitizedTagName := utils.SanitizeText(tagName, utils.TEXT)
	sanitizedTagColor := utils.SanitizeText(tagColor, utils.PASSWORD)

	// If sanitization removed all characters, reject the input
	if sanitizedTagName == "" || sanitizedTagColor == "" {
		return 0, errors.New("tag name or color contains only invalid characters")
	}

	if len(sanitizedTagColor) > 0 && sanitizedTagColor[0] != '#' {
		sanitizedTagColor = "#" + sanitizedTagColor
	}

	// Begin transaction to ensure data consistency
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Get the user ID from the session
	var userID int
	err = tx.QueryRow("SELECT userid FROM SESSION WHERE id = ?", sessionID).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("no session exists for the provided sessionID")
		}
		return 0, err
	}

	// Check if the project exists
	var existingProjectID int
	err = tx.QueryRow("SELECT id FROM PROJECT WHERE id = ?", projectID).Scan(&existingProjectID)
	if err == sql.ErrNoRows {
		return 0, errors.New("project does not exist")
	} else if err != nil {
		return 0, err
	}

	// Check if the user has write privileges on the project
	var hasWritePrivileges bool
	err = tx.QueryRow(`
	SELECT 1
	FROM PROJECT_MEMBER pm
	JOIN PROJECT_ROLE pr ON pm.projectid = pr.projectid
	JOIN ORG_ROLE orl ON pr.id = orl.id
	WHERE pm.userid = ? AND pm.projectid = ? AND orl.can_write = 1
	`, userID, projectID).Scan(&hasWritePrivileges)
	if err != nil {
		return 0, err
	}
	if !hasWritePrivileges {
		return 0, errors.New("user does not have write privileges on this project")
	}

	// Check if the tag name already exists in the project
	var existingTagID int
	err = tx.QueryRow("SELECT id FROM TAG WHERE projectid = ? AND name = ?", projectID, sanitizedTagName).Scan(&existingTagID)
	if err == nil {
		return 0, errors.New("tag name already exists in the project")
	}
	if err != sql.ErrNoRows {
		return 0, err
	}

	// Insert new tag and get its ID
	result, err := tx.Exec("INSERT INTO TAG(projectid, name, color) VALUES(?, ?, ?)", projectID, sanitizedTagName, sanitizedTagColor)
	if err != nil {
		return 0, err
	}

	tagID64, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	tagID := int(tagID64)

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return 0, err
	}

	// Return the ID of the newly created tag
	return tagID, nil
}
