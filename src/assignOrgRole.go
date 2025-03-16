package backend

import (
	"database/sql"

	_ "modernc.org/sqlite" // SQLite driver for database/sql
)

// User represents a user in the system.
type user struct {
	ID        int
	Validated bool
	Roles     map[int]bool
}

// Organization represents an organization with users and executive roles.
type Organization struct {
	ID    int
	Users map[int]*user
	Execs map[int]bool
}

// getUserByID retrieves a user from the database.
func getUserByID(db *sql.DB, userID int) (*user, error) {
	user := &user{ID: userID, Roles: make(map[int]bool)}

	err := db.QueryRow("SELECT (verifyid IS NOT NULL) FROM USER WHERE id = ?", userID).Scan(&user.Validated)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT roleid FROM ORG_MEMBER_ROLE WHERE memberid = (SELECT id FROM ORG_MEMBER WHERE userid = ?)", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var roleID int
		if err := rows.Scan(&roleID); err != nil {
			return nil, err
		}
		user.Roles[roleID] = true
	}

	return user, nil
}

// getOrganizationByID retrieves an organization from the database.
func getOrganizationByID(db *sql.DB, orgID int) (*Organization, error) {
	org := &Organization{ID: orgID, Users: make(map[int]*user), Execs: make(map[int]bool)}

	err := db.QueryRow("SELECT id FROM ORGANIZATION WHERE id = ?", orgID).Scan(&org.ID)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT userid FROM ORG_MEMBER WHERE orgid = ?", orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		user, err := getUserByID(db, userID)
		if err != nil {
			return nil, err
		}
		org.Users[userID] = user
	}

	execRows, err := db.Query("SELECT userid FROM ORG_MEMBER_ROLE WHERE roleid = (SELECT id FROM ORG_ROLE WHERE orgid = ? AND can_exec = 1)", orgID)
	if err != nil {
		return nil, err
	}
	defer execRows.Close()

	for execRows.Next() {
		var userID int
		if err := execRows.Scan(&userID); err != nil {
			return nil, err
		}
		org.Execs[userID] = true
	}

	return org, nil
}

// assignOrgRole promotes User B to a role within an organization.
// User A (acting user) is referenced by sessionID.
// User B (target user) is referenced by userID.
func assignOrgRole(db *sql.DB, sessionID, userID, orgID, newRoleID int) error {
	// Get User A's member ID and role in the organization
	var sessionMemberID, sessionRoleID int
	err := db.QueryRow("SELECT id FROM ORG_MEMBER WHERE userid = ? AND orgid = ?", sessionID, orgID).Scan(&sessionMemberID)
	if err != nil {
		return err
	}

	err = db.QueryRow("SELECT roleid FROM ORG_MEMBER_ROLE WHERE memberid = ?", sessionMemberID).Scan(&sessionRoleID)
	if err != nil {
		return err
	}

	// Check if User A has permission to promote (can_exec = 1)
	var canExec bool
	err = db.QueryRow("SELECT can_exec FROM ORG_ROLE WHERE id = ? AND orgid = ?", sessionRoleID, orgID).Scan(&canExec)
	if err != nil {
		return err
	}
	if !canExec {
		return err
	}

	// Verify that User B is a member of the organization
	var memberID int
	err = db.QueryRow("SELECT id FROM ORG_MEMBER WHERE userid = ? AND orgid = ?", userID, orgID).Scan(&memberID)
	if err != nil {
		return err
	}

	// Ensure User B is validated
	var isValidated bool
	err = db.QueryRow("SELECT (verifyid IS NOT NULL) FROM USER WHERE id = ?", userID).Scan(&isValidated)
	if err != nil {
		return err
	}
	if !isValidated {
		return err
	}

	// Assign the new role to User B
	_, err = db.Exec("UPDATE ORG_MEMBER_ROLE SET roleid = ? WHERE memberid = ?", newRoleID, memberID)
	return err
}
