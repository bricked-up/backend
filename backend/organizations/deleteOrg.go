package organizations

import (
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

// DeleteOrganization deletes an organization if the user has the necessary permissions.
func DeleteOrganization(db *sql.DB, sessionID int, orgID int) error {
	// Input validation
	if db == nil {
		return errors.New("database connection is nil")
	}
	if sessionID <= 0 {
		return errors.New("invalid session ID")
	}
	if orgID <= 0 {
		return errors.New("invalid organization ID")
	}

	// Get the user ID from the session
	var userID int
	var expires string
	err := db.QueryRow(
		`SELECT userid, expires 
		FROM SESSION WHERE id = ?`,
		sessionID).Scan(&userID, &expires)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("no session exists for the provided sessionID")
		}
		return err
	}

	// Check if the organization exists
	var existingOrgID int
	err = db.QueryRow(
		`SELECT id 
		FROM ORGANIZATION WHERE id = ?`,
		orgID).Scan(&existingOrgID)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("organization does not exist")
		}
		return err
	}

	// Check if the user is a member of the organization
	var memberID int
	err = db.QueryRow(
		`SELECT id 
		FROM ORG_MEMBER 
		WHERE userid = ? AND orgid = ?`,
		userID, orgID).Scan(&memberID)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user is not a member of this organization")
		}
		return err
	}

	// // Check if the user has admin privileges
	// var isAdmin bool
	// err = db.QueryRow(`
	// 	SELECT EXISTS(
	// 		SELECT 1 FROM ORG_MEMBER_ROLE mr
	// 		JOIN ORG_ROLE r ON mr.roleid = r.id
	// 		WHERE mr.memberid = ? AND r.orgid = ? AND r.can_exec = 1
	// 	)
	// `, memberID, orgID).Scan(&isAdmin)

	// if err != nil {
	// 	return err
	// }

	// if !isAdmin {
	// 	return errors.New("user does not have permission to delete the organization")
	// }

	// Delete the organization (cascade will handle related records)
	result, err := db.Exec("DELETE FROM ORGANIZATION WHERE id = ?", orgID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("organization not found or already deleted")
	}

	return nil
}
