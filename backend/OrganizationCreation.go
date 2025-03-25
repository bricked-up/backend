package backend

import (
	"brickedup/backend/utils"
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

// CreateOrganization creates a new organization and assigns the user (from the session) to it as an admin.
// It takes sessionID (int) and orgName (string) as parameters.
func CreateOrganization(db *sql.DB, sessionID int, orgName string) (int, error) {
	// Check if orgName is provided
	if orgName == "" {
		return 0, errors.New("missing orgName")
	}

	// In the function
	sanitizedOrgName := utils.SanitizeText(orgName, utils.TEXT)

	// If sanitization removed all characters, reject the input
	if sanitizedOrgName == "" {
		return 0, errors.New("organization name contains only invalid characters")
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

	// Check if the organization name already exists
	var existingOrgID int
	err = tx.QueryRow("SELECT id FROM ORGANIZATION WHERE name = ?", sanitizedOrgName).Scan(&existingOrgID)
	if err == nil {
		return 0, errors.New("organization name already exists")
	}
	if err != sql.ErrNoRows {
		return 0, err
	}

	// Insert new organization and get its ID
	result, err := tx.Exec("INSERT INTO ORGANIZATION(name) VALUES(?)", sanitizedOrgName)
	if err != nil {
		return 0, err
	}

	orgID64, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	orgID := int(orgID64)

	// Create admin role for the organization with rwx permissions
	roleResult, err := tx.Exec(`
		INSERT INTO ORG_ROLE (orgid, name, can_read, can_write, can_exec)
		VALUES (?, 'admin', 1, 1, 1)`, orgID)
	if err != nil {
		return 0, err
	}

	roleID64, err := roleResult.LastInsertId()
	if err != nil {
		return 0, err
	}
	roleID := int(roleID64)

	// Insert into ORG_MEMBER
	memberResult, err := tx.Exec("INSERT INTO ORG_MEMBER (userid, orgid) VALUES (?, ?)", userID, orgID)
	if err != nil {
		return 0, err
	}

	memberID64, err := memberResult.LastInsertId()
	if err != nil {
		return 0, err
	}
	memberID := int(memberID64)

	// Assign role to member
	_, err = tx.Exec("INSERT INTO ORG_MEMBER_ROLE (memberid, roleid) VALUES (?, ?)", memberID, roleID)
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
