package backend

import (
	"database/sql"
	"errors"
	"strconv"

	_ "modernc.org/sqlite"
)

// CreateOrganization creates a new organization and assigns the user (from the session) to it as an admin.
// It takes sessionID (int) and orgName (string) as parameters.
func CreateOrganization(db *sql.DB, sessionID int, orgName string) (int, error) {
	// Check if orgName is provided
	if orgName == "" {
		return 0, errors.New("missing orgName")
	}

	// Get the user ID from the session
	var userID int
	var err error
	err = db.QueryRow("SELECT userid FROM SESSION WHERE id = ?", sessionID).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("no session found for session ID " + strconv.Itoa(sessionID))
		}
		return 0, err
	}

	// Check if the organization name already exists
	var existingOrgID int
	err = db.QueryRow("SELECT id FROM ORGANIZATION WHERE name = ?", orgName).Scan(&existingOrgID)
	if err == nil {
		return 0, errors.New("organization name already exists")
	}
	if err != sql.ErrNoRows {
		return 0, err
	}

	// Insert new organization and get its ID
	result, err := db.Exec("INSERT INTO ORGANIZATION(name) VALUES(?)", orgName)
	if err != nil {
		return 0, err
	}

	orgID64, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	orgID := int(orgID64)

	// Create admin role for the organization with rwx permissions
	roleResult, err := db.Exec(`
		INSERT INTO ORG_ROLE (orgid, name, can_read, can_write, can_exec)
		VALUES (?, 'admin', 1, 1, 1)`, orgID)
	if err != nil {
		return 0, err
	}

	// Assign user to the organization as the admin (role_id = 1 for admin)
	_, err = db.Exec("INSERT INTO organization_members (user_id, organization_id, role_id) VALUES (?, ?, ?)", userID, orgID, 1)
	if err != nil {
		return 0, err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return 0, err
	}

	// Return the ID of the newly created organization
	return orgID, nil
}
