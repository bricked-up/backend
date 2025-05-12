package projects

import (
	"brickedup/backend/utils"
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

// CreateProj creates a new project and assigns the user (from the session) to it as an admin.
func CreateProj(
	db *sql.DB, 
	sessionID int,
	orgid int,
	name string,
	budget int, 
	charter string) error {

	if name == "" {
		return errors.New("missing project name")
	}

	// In the function
	name = utils.SanitizeText(name, utils.TEXT)
	charter = utils.SanitizeText(charter, utils.TEXT)

	// If sanitization removed all characters, reject the input
	if name == "" {
		return errors.New("project name contains only invalid characters")
	}

	// Begin transaction to ensure data consistency
	tx, err := db.Begin()
	if err != nil {
		return err
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
			return errors.New("no session exists for the provided sessionID")
		}
		return err
	}

	result, err := tx.Exec(
		"INSERT INTO PROJECT(name, budget, charter, orgid, archived) VALUES(?, ?, ?, ?, 0)", 
		name, budget, charter, orgid)

	if err != nil {
		return err
	}

	projID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	roleResult, err := tx.Exec(`
		INSERT INTO PROJECT_ROLE (projectid, name, can_read, can_write, can_exec)
		VALUES (?, 'admin', 1, 1, 1)`, projID)

	if err != nil {
		return err
	}

	roleID64, err := roleResult.LastInsertId()
	if err != nil {
		return err
	}
	roleID := int(roleID64)

	// Insert into ORG_MEMBER
	memberResult, err := tx.Exec(
		`INSERT INTO PROJECT_MEMBER (userid, projectid) 
		VALUES (?, ?)`, 
		userID, projID)

	if err != nil {
		return err
	}

	memberID64, err := memberResult.LastInsertId()
	if err != nil {
		return err
	}
	memberID := int(memberID64)

	// Assign role to member
	_, err = tx.Exec(
		`INSERT INTO PROJECT_MEMBER_ROLE (memberid, roleid) 
		VALUES (?, ?)`,
		memberID, roleID)

	if err != nil {
		return err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return err
	}

	// Return the ID of the newly created project
	return nil
}
