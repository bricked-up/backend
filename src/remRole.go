package backend

import (
	"database/sql"
	"errors"
)

// removeUserRole removes the role of user B in a project
func removeUserRole(db *sql.DB, sessionID, userBID string, roleID, projectID int) error {
	var userAID int
	var userAValidated bool
	var userBExists bool
	var userBValidated bool
	var userBInProject bool
	var userBHasRole bool
	var userAPermissions bool

	// Get the user A's ID from the session (assuming sessionID maps to user A)
	err := db.QueryRow(`SELECT userid FROM SESSION WHERE id = ?`, sessionID).Scan(&userAID)
	if err != nil || userAID == 0 {
		return errors.New("invalid session ID")
	}

	// Check if user A is validated
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM VERIFY_USER WHERE id = (SELECT verifyid FROM USER WHERE id = ?))`, userAID).Scan(&userAValidated)
	if err != nil || !userAValidated {
		return errors.New("user A is not validated")
	}

	// Check if user B exists
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM USER WHERE id = ?)`, userBID).Scan(&userBExists)
	if err != nil || !userBExists {
		return errors.New("user B does not exist")
	}

	// Check if user B is validated
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM VERIFY_USER WHERE id = (SELECT verifyid FROM USER WHERE id = ?))`, userBID).Scan(&userBValidated)
	if err != nil || !userBValidated {
		return errors.New("user B is not validated")
	}

	// Check if user B is part of the project
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM PROJECT_MEMBER WHERE userid = ? AND projectid = ?)`, userBID, projectID).Scan(&userBInProject)
	if err != nil || !userBInProject {
		return errors.New("user B is not part of the project")
	}

	// Check if user B has the specified role in the project
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM PROJECT_MEMBER_ROLE WHERE memberid = (SELECT id FROM PROJECT_MEMBER WHERE userid = ? AND projectid = ?) AND roleid = ?)`, userBID, projectID, roleID).Scan(&userBHasRole)
	if err != nil || !userBHasRole {
		return errors.New("user B does not have the specified role")
	}

	// Check if user A has exec permissions in the project
	err = db.QueryRow(`SELECT can_exec FROM PROJECT_MEMBER_ROLE WHERE memberid = (SELECT id FROM PROJECT_MEMBER WHERE userid = ? AND projectid = ?)`, userAID, projectID).Scan(&userAPermissions)
	if err != nil || !userAPermissions {
		return errors.New("user A lacks exec permissions")
	}

	// Remove user B's role in the project
	_, err = db.Exec(`DELETE FROM PROJECT_MEMBER_ROLE WHERE memberid = (SELECT id FROM PROJECT_MEMBER WHERE userid = ? AND projectid = ?) AND roleid = ?`, userBID, projectID, roleID)
	if err != nil {
		return err
	}

	return nil
}
