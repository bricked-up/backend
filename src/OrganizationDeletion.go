package backend

import (
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

func DeleteOrganization(sessionID, orgID string) error {
	// Check if sessionID and orgID are provided
	if sessionID == "" || orgID == "" {
		return errors.New("missing sessionID or orgID")
	}

	// Open the database
	db, err := sql.Open("sqlite", "backend/sql/init.sql")
	if err != nil {
		return err
	}
	defer db.Close()

	// Validate session ID and retrieve the user ID associated with it
	var userID int
	err = db.QueryRow("SELECT user_id FROM sessions WHERE session_id = ?", sessionID).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("no session found for session ID " + sessionID)
		}
		return err
	}

	// Check if the organization exists
	var existingOrgID int
	err = db.QueryRow("SELECT id FROM organizations WHERE id = ?", orgID).Scan(&existingOrgID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("organization does not exist")
		}
		return err
	}

	// Retrieve User ID from organization_roles
	var roleID int
	err = db.QueryRow("SELECT role_id FROM organization_roles WHERE organization_id = ? AND user_id = ?", orgID, userID).Scan(&roleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user does not have permission to delete the organization")
		}
		return err
	}

	// Check if the role has exec permission
	var canExec bool
	err = db.QueryRow("SELECT can_exec FROM roles WHERE id = ?", roleID).Scan(&canExec)
	if err != nil {
		return err
	}
	if !canExec {
		return errors.New("user does not have permission to delete the organization")
	}

	// Delete the organization
	_, err = db.Exec("DELETE FROM organizations WHERE id = ?", orgID)
	if err != nil {
		return err
	}

	return nil
}
