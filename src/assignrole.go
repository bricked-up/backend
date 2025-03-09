package backend

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

/*AssignProjectRoleToUser assigns a new role to a user(assignee) in a project.
It ensures that assignor has exec permission, assignee is validated and the role is not already assigned.*/

func assignProjectRoleToUser(db *sql.DB, sessionIDA, sessionIDB string, roleID, projectID int) error {
	//Validate that the project exists
	var projectExists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM PROJECT WHERE id = ?)", projectID).Scan(&projectExists)
	if err != nil || !projectExists {
		return err
	}

	//Validate Assignor's session and permissions
	_, err = validateSessionAndExecPermission(db, sessionIDA, projectID)
	if err != nil {
		return err
	}

	//Validate User B's session and get their UserID
	userBID, err := validateUserSessionMembershipAndVerification(db, sessionIDB, projectID)
	if err != nil {
		return err
	}

	//Check if Assignee already has the role
	var hasRole bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM PROJECT_MEMBER WHERE userid = ? AND projectid = ? AND roleid = ?)", userBID, projectID, roleID).Scan(&hasRole)
	if err != nil || hasRole {
		return err
	}

	//Update Assignee's role in the project
	_, err = db.Exec("UPDATE PROJECT_MEMBER SET roleid = ? WHERE userid = ? AND projectid = ?", roleID, userBID, projectID)
	return err
}

// ValidateSessionAndExecPermission ensures Assignor's session is valid and they have exec permission in the project.
func validateSessionAndExecPermission(db *sql.DB, sessionID string, projectID int) (int, error) {
	var userID int
	var hasExecPermission bool

	//Validate Assignor's session
	err := db.QueryRow("SELECT userid FROM SESSION WHERE token = ? AND expires > ?", sessionID, time.Now()).Scan(&userID)
	if err != nil {
		return 0, err
	}

	//Check if assignor has exec permission in the project
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM PROJECT_MEMBER_ROLE pmr JOIN PROJECT_ROLE pr ON pmr.roleid = pr.id JOIN PROJECT_MEMBER pm ON pmr.memberid = pm.id WHERE pm.userid = ? AND pm.projectid = ? AND pr.can_exec = TRUE)", userID, projectID).Scan(&hasExecPermission)
	if err != nil {
		return 0, err
	}
	if !hasExecPermission {
		return 0, fmt.Errorf("user does not have exec permission")
	}

	return userID, nil
}

// ValidateUserSessionMembershipAndVerification ensures assignee's session is valid, they are a project member and verified
func validateUserSessionMembershipAndVerification(db *sql.DB, sessionID string, projectID int) (int, error) {
	var userID int
	var isMember bool
	var isVerified int

	//Validate Assignee's session
	err := db.QueryRow("SELECT userid FROM SESSION WHERE token = ? AND expires > ?", sessionID, time.Now()).Scan(&userID)
	if err != nil {
		return 0, err
	}

	//Check if Assignee is a member of the project
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM PROJECT_MEMBER WHERE userid = ? AND projectid = ?)", userID, projectID).Scan(&isMember)
	if err != nil {
		return 0, err
	}
	if !isMember {
		return 0, err
	}

	//Check if assignee is verified
	err = db.QueryRow("SELECT verifyid FROM USER WHERE id = ?", userID).Scan(&isVerified)
	if err != nil {
		return 0, err
	}
	if isVerified == 0 {
		return 0, err
	}

	return userID, nil
}
