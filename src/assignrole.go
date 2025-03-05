package backend

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

var db *sql.DB

/*assignProjectRoleToUser assigns a new role to User B in a project.
It ensures that User A has exec permission, User B is validated and the role is not already assigned.*/

func assignProjectRoleToUser(sessionIDA, sessionIDB string, roleID, projectID int) error {
	// Validate that the project exists
	var projectExists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM PROJECT WHERE id = ?)", projectID).Scan(&projectExists)
	if err != nil {
		return fmt.Errorf("failed to check project existence: %v", err)
	}
	if !projectExists {
		return fmt.Errorf("project %d does not exist", projectID)
	}

	//Validate User A's session and permissions
	_, err = validateSessionAndExecPermission(sessionIDA, projectID)
	if err != nil {
		return fmt.Errorf("assignor validation failed: %v", err)
	}

	//Validate User B's session, membership and verification status
	userBID, err := validateUserSessionMembershipAndVerification(sessionIDB, projectID)
	if err != nil {
		return fmt.Errorf("assignee validation failed: %v", err)
	}

	//Check if User B already has the role
	var hasRole bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM PROJECT_MEMBER WHERE userid = ? AND projectid = ? AND roleid = ?)", userBID, projectID, roleID).Scan(&hasRole)
	if err != nil {
		return fmt.Errorf("failed to check if assignee already has the role: %v", err)
	}
	if hasRole {
		return fmt.Errorf("assignee already has the role %d in project %d", roleID, projectID)
	}

	//Update User B's role in the project
	_, err = db.Exec("UPDATE PROJECT_MEMBER SET roleid = ? WHERE userid = ? AND projectid = ?", roleID, userBID, projectID)
	if err != nil {
		return fmt.Errorf("failed to update assignee's role: %v", err)
	}

	return nil
}

// validateSessionAndExecPermission ensures User A's session is valid and they have exec permission in the project.
func validateSessionAndExecPermission(sessionID string, projectID int) (int, error) {
	var userID int
	var execPermission bool

	//Validate User A's session
	err := db.QueryRow("SELECT userid FROM SESSION WHERE token = ? AND expires > ?", sessionID, time.Now()).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("Invalid session for assignor: %v", err)
	}

	//Check if User A has exec permission in the project
	err = db.QueryRow("SELECT exec FROM PROJECT_MEMBER WHERE userid = ? AND projectid = ?", userID, projectID).Scan(&execPermission)
	if err != nil {
		return 0, fmt.Errorf("Failed to check User A's permissions: %v", err)
	}
	if !execPermission {
		return 0, fmt.Errorf("Assignor does not have exec permission in project %d", projectID)
	}

	return userID, nil
}

// validateUserSessionMembershipAndVerification ensures User B's session is valid, they are a project member, and verified
func validateUserSessionMembershipAndVerification(sessionID string, projectID int) (int, error) {
	var userID int
	var isMember bool
	var isVerified int

	//Validate User B's session
	err := db.QueryRow("SELECT userid FROM SESSION WHERE token = ? AND expires > ?", sessionID, time.Now()).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("Invalid session for User B: %v", err)
	}

	//Check if User B is a member of the project
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM PROJECT_MEMBER WHERE userid = ? AND projectid = ?)", userID, projectID).Scan(&isMember)
	if err != nil {
		return 0, fmt.Errorf("Failed to check User B's membership: %v", err)
	}
	if !isMember {
		return 0, fmt.Errorf("User B is not a member of project %d", projectID)
	}

	//Check if User B is verified
	err = db.QueryRow("SELECT verifyid FROM USER WHERE id = ?", userID).Scan(&isVerified)
	if err != nil {
		return 0, fmt.Errorf("Failed to check User B's verification status: %v", err)
	}
	if isVerified == 0 {
		return 0, fmt.Errorf("User B is not verified")
	}

	return userID, nil
}
