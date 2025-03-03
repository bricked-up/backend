package backend

import (
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

const dbPath = "backend/sql/BrickedUpDatabase.sql"

// CreateOrganization creates a new organization and assigns the user (from the session) to it as an admin.
// It takes sessionID and orgName as parameters instead of extracting them from the request.
func CreateOrganization(sessionID, orgName string) (int, error) {
	// Check if sessionID and orgName are provided
	if sessionID == "" || orgName == "" {
		return 0, errors.New("missing sessionID or orgName")
	}

	// Open the database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	// Validate session ID and retrieve the user ID associated with it
	var userID int
	err = db.QueryRow("SELECT user_id FROM sessions WHERE session_id = ?", sessionID).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("no session found for session ID " + sessionID)
		}
		return 0, err
	}

	// Check if the organization name already exists
	var existingOrgID int
	err = db.QueryRow("SELECT id FROM organizations WHERE name = ?", orgName).Scan(&existingOrgID)
	if err == nil {
		return 0, errors.New("organization name already exists")
	}
	if err != sql.ErrNoRows {
		return 0, err
	}

	// Insert new organization and get its ID
	var orgID int
	err = db.QueryRow("INSERT INTO organizations(name) VALUES(?) RETURNING id", orgName).Scan(&orgID)
	if err != nil {
		return 0, err
	}

	// Create admin role for the organization with rwx permissions
	_, err = db.Exec(`
		INSERT INTO organization_roles (organization_id, role_name, can_read, can_write, can_execute)
		VALUES (?, 'admin', 1, 1, 1)`, orgID)
	if err != nil {
		return 0, err
	}

	// Assign user to the organization as the admin (role_id = 1 for admin)
	_, err = db.Exec("INSERT INTO organization_members (user_id, organization_id, role_id) VALUES (?, ?, ?)", userID, orgID, 1)
	if err != nil {
		return 0, err
	}

	// Return the ID of the newly created organization
	return orgID, nil
}
