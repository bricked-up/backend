package backend

import (
	"fmt"
	"log"
	"time"
)

//assignOrgRoleToUser assigns a role to a user in an organization.

//   - sessionIDA: Session ID of User A (the user making the change).
//   - sessionIDB: Session ID of User B (the user whose role is being changed).
//   - roleID: The ID of the new role to assign to User B.
//   - orgID: The ID of the organization where the role change is happening
//
// It returns an error if the assignment fails.
func assignOrgRoleToUser(sessionIDA, sessionIDB string, roleID, orgID int) error {
	//Validate User A's session and permissions
	userAID, err := validateSessionAndPermissions(sessionIDA, orgID)
	if err != nil {
		return fmt.Errorf("user A validation failed: %v", err)
	}

	//Validate User B's session and membership in the organization
	userBID, err := validateUserSessionAndMembership(sessionIDB, orgID)
	if err != nil {
		return fmt.Errorf("user B validation failed: %v", err)
	}

	//Check if User B already has the role
	var hasRole bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM ORG_MEMBER WHERE userid = ? AND orgid = ? AND roleid = ?)", userBID, orgID, roleID).Scan(&hasRole)
	if err != nil {
		return fmt.Errorf("failed to check if user B already has the role: %v", err)
	}
	if hasRole {
		return fmt.Errorf("user B already has the role %d in organization %d", roleID, orgID)
	}

	//Update User B's role in the organization
	_, err = db.Exec("UPDATE ORG_MEMBER SET roleid = ? WHERE userid = ? AND orgid = ?", roleID, userBID, orgID)
	if err != nil {
		return fmt.Errorf("failed to update user B's role: %v", err)
	}

	//Log the role change
	logRoleChange(userAID, userBID, roleID, orgID)

	return nil
}

// validateSessionAndPermissions validates User A's session and checks if they have exec permission in the organization.
func validateSessionAndPermissions(sessionID string, orgID int) (int, error) {
	var userID int
	var execPermission bool

	//Validate User A's session
	err := db.QueryRow("SELECT userid FROM SESSION WHERE token = ? AND expires > ?", sessionID, time.Now()).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("invalid session for User A: %v", err)
	}

	//Check if User A has exec permission in the organization
	err = db.QueryRow("SELECT exec FROM ORG_MEMBER WHERE userid = ? AND orgid = ?", userID, orgID).Scan(&execPermission)
	if err != nil {
		return 0, fmt.Errorf("failed to check User A's permissions: %v", err)
	}
	if !execPermission {
		return 0, fmt.Errorf("user A does not have exec permission in organization %d", orgID)
	}

	return userID, nil
}

// validateUserSessionAndMembership validates User B's session and checks if they are a member of the organization
func validateUserSessionAndMembership(sessionID string, orgID int) (int, error) {
	var userID int

	//Validate User B's session
	err := db.QueryRow("SELECT userid FROM SESSION WHERE token = ? AND expires > ?", sessionID, time.Now()).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("invalid session for User B: %v", err)
	}

	//Check if User B is a member of the organization
	var isMember bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM ORG_MEMBER WHERE userid = ? AND orgid = ?)", userID, orgID).Scan(&isMember)
	if err != nil {
		return 0, fmt.Errorf("failed to check User B's membership: %v", err)
	}
	if !isMember {
		return 0, fmt.Errorf("user B is not a member of organization %d", orgID)
	}

	return userID, nil
}

// logRoleChange logs the role change for auditing purposes.
func logRoleChange(userAID, userBID, roleID, orgID int) {
	//Insert a log entry into the database
	_, err := db.Exec("INSERT INTO ROLE_CHANGE_LOG (user_a_id, user_b_id, role_id, org_id, change_time) VALUES (?, ?, ?, ?, ?)", userAID, userBID, roleID, orgID, time.Now())
	if err != nil {
		log.Printf("Failed to log role change: %v", err)
	}
}
