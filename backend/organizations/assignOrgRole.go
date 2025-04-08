package organizations

import (
	"database/sql"

	_ "modernc.org/sqlite" // SQLite driver for database/sql
)

// AssignOrgRole promotes User B to a role within an organization.
// User A (acting user) is referenced by sessionID.
// User B (target user) is referenced by userID.
func AssignOrgRole(db *sql.DB, sessionID, userID, orgID, newRoleID int) error {
	// Get User A's member ID and role in the organization
	var sessionMemberID, sessionRoleID int
	err := db.QueryRow(
		`SELECT id 
		FROM ORG_MEMBER 
		WHERE userid = ? AND orgid = ?`, 
		sessionID, orgID).Scan(&sessionMemberID)

	if err != nil {
		return err
	}

	err = db.QueryRow(
		`SELECT roleid 
		FROM ORG_MEMBER_ROLE 
		WHERE memberid = ?`, 
		sessionMemberID).Scan(&sessionRoleID)

	if err != nil {
		return err
	}

	// Check if User A has permission to promote (can_exec = 1)
	var canExec bool
	err = db.QueryRow(
		`SELECT can_exec 
		FROM ORG_ROLE WHERE id = ? AND orgid = ?`, 
		sessionRoleID, orgID).Scan(&canExec)

	if err != nil {
		return err
	}
	if !canExec {
		return err
	}

	// Verify that User B is a member of the organization
	var memberID int
	err = db.QueryRow(
		`SELECT id 
		FROM ORG_MEMBER 
		WHERE userid = ? AND orgid = ?`, 
		userID, orgID).Scan(&memberID)

	if err != nil {
		return err
	}

	// Ensure User B is validated
	var isValidated bool
	err = db.QueryRow(
		`SELECT (verifyid IS NOT NULL) 
		FROM USER 
		WHERE id = ?`, 
		userID).Scan(&isValidated)

	if err != nil {
		return err
	}
	if !isValidated {
		return err
	}

	// Assign the new role to User B
	_, err = db.Exec(
		`UPDATE ORG_MEMBER_ROLE 
		SET roleid = ? 
		WHERE memberid = ?`, 
		newRoleID, memberID)

	return err
}
