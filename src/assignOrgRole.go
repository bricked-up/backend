package backend

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // SQLite driver for database/sql
)

// User and Organization structs are created to simplify importing database values in functions

type User struct {
	ID        string
	Validated bool
	Roles     map[int]bool
}

type Organization struct {
	ID    int
	Users map[string]*User
	Execs map[string]bool
}

// getUserByID retrieves a user from the database.
func getUserByID(db *sql.DB, userID string) (*User, error) {
	user := &User{Roles: make(map[int]bool)}

	err := db.QueryRow("SELECT id, verified FROM USER WHERE id = ?", userID).Scan(&user.ID, &user.Validated)
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
	org := &Organization{ID: orgID, Users: make(map[string]*User), Execs: make(map[string]bool)}

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
		var userID string
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
		var userID string
		if err := execRows.Scan(&userID); err != nil {
			return nil, err
		}
		org.Execs[userID] = true
	}

	return org, nil
}

// getUserBySession retrieves the user ID from a session ID.
func getUserBySession(db *sql.DB, sessionID string) (string, error) {
	var userID string
	err := db.QueryRow("SELECT userid FROM SESSION WHERE id = ?", sessionID).Scan(&userID)
	if err != nil {
		return err.Error(), err
	}
	return userID, nil
}

// assignOrgRole promotes user B to a role within an organization.
func assignOrgRole(db *sql.DB, userEmail, roleName string, orgID, newRoleID int) error {
	// Fetch the user ID from the USER table
	var userID int
	err := db.QueryRow("SELECT id FROM USER WHERE email = ?", userEmail).Scan(&userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %v", err)
	}
	fmt.Printf("Found user: %d (email: %s)\n", userID, userEmail) // Debug info

	// Fetch the organization details
	var orgName string
	err = db.QueryRow("SELECT name FROM ORGANIZATION WHERE id = ?", orgID).Scan(&orgName)
	if err != nil {
		return fmt.Errorf("failed to find organization: %v", err)
	}
	fmt.Printf("Found organization: %d (name: %s)\n", orgID, orgName) // Debug info

	// Check if the user is a member of the organization
	var memberID int
	err = db.QueryRow("SELECT id FROM ORG_MEMBER WHERE userid = ? AND orgid = ?", userID, orgID).Scan(&memberID)
	if err != nil {
		return fmt.Errorf("user is not part of the organization: %v", err)
	}
	fmt.Printf("Found organization member: %d\n", memberID) // Debug info

	// Fetch the current role of the user in the organization
	var currentRoleID int
	err = db.QueryRow(`
		SELECT roleid 
		FROM ORG_MEMBER_ROLE 
		WHERE memberid = (SELECT id FROM ORG_MEMBER WHERE userid = ? AND orgid = ?)`,
		userID, orgID).Scan(&currentRoleID)
	if err != nil {
		return fmt.Errorf("failed to find current role: %v", err)
	}
	fmt.Printf("Current role ID: %d\n", currentRoleID) // Debug info

	// Check if the user has the necessary permissions to promote
	var canExec bool
	err = db.QueryRow(`
		SELECT can_exec 
		FROM ORG_ROLE 
		WHERE id = ? AND orgid = ?`,
		currentRoleID, orgID).Scan(&canExec)
	if err != nil {
		return fmt.Errorf("failed to check permissions: %v", err)
	}
	fmt.Printf("User has can_exec permission: %v\n", canExec) // Debug info
	if !canExec {
		return fmt.Errorf("no exec permissions for promotion")
	}

	// Check if the user is validated
	var sessionID int
	err = db.QueryRow("SELECT id FROM SESSION WHERE userid = ?", userID).Scan(&sessionID)
	if err != nil {
		return fmt.Errorf("attempted to change user role, user not validated")
	}
	fmt.Printf("User session ID: %d\n", sessionID) // Debug info

	// Assign the new role to the user in the organization
	_, err = db.Exec(`
		UPDATE ORG_MEMBER_ROLE 
		SET roleid = ? 
		WHERE memberid = ? AND roleid = ?`,
		newRoleID, memberID, currentRoleID)
	if err != nil {
		return fmt.Errorf("failed to assign new role: %v", err)
	}

	return nil
}
