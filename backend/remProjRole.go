package backend

import (
	"database/sql"
	"errors"
)

// removeUserRole removes the role of user B in a project.
// sessionID is used to authenticate user A (the initiator).
func removeUserRole(db *sql.DB, sessionID int, userBID string, roleID, projectID int) error {
	var userAID int
	var isUserAValidated, isUserBValidated bool
	var userBExists, userBInProject, userBHasRole, userAHasExec bool

	// Get user A from session
	err := db.QueryRow(`SELECT userid FROM SESSION WHERE id = ?`, sessionID).Scan(&userAID)
	if err != nil {
		return errors.New("invalid session ID")
	}

	// Check if user A is verified
	err = db.QueryRow(`SELECT verified FROM USER WHERE id = ?`, userAID).Scan(&isUserAValidated)
	if err != nil || !isUserAValidated {
		return errors.New("user A is not validated")
	}

	// Check if user B exists
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM USER WHERE id = ?)`, userBID).Scan(&userBExists)
	if err != nil || !userBExists {
		return errors.New("user B does not exist")
	}

	// Check if user B is verified
	err = db.QueryRow(`SELECT verified FROM USER WHERE id = ?`, userBID).Scan(&isUserBValidated)
	if err != nil || !isUserBValidated {
		return errors.New("user B is not validated")
	}

	// Check if user B is in the project
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM PROJECT_MEMBER WHERE userid = ? AND projectid = ?)`, userBID, projectID).Scan(&userBInProject)
	if err != nil || !userBInProject {
		return errors.New("user B is not part of the project")
	}

	// Check if user B has the role
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM PROJECT_MEMBER_ROLE
			WHERE memberid = (
				SELECT id FROM PROJECT_MEMBER WHERE userid = ? AND projectid = ?
			) AND roleid = ?
		)`, userBID, projectID, roleID).Scan(&userBHasRole)
	if err != nil || !userBHasRole {
		return errors.New("user B does not have the specified role")
	}

	// Check if user A has exec permission in the project
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM PROJECT_MEMBER_ROLE pmr
			JOIN PROJECT_ROLE pr ON pmr.roleid = pr.id
			WHERE pmr.memberid = (
				SELECT id FROM PROJECT_MEMBER WHERE userid = ? AND projectid = ?
			) AND pr.can_exec = 1
		)`, userAID, projectID).Scan(&userAHasExec)
	if err != nil || !userAHasExec {
		return errors.New("user A lacks exec permissions")
	}

	// Delete the role from user B
	_, err = db.Exec(`
		DELETE FROM PROJECT_MEMBER_ROLE
		WHERE memberid = (
			SELECT id FROM PROJECT_MEMBER WHERE userid = ? AND projectid = ?
		) AND roleid = ?
	`, userBID, projectID, roleID)
	if err != nil {
		return err
	}

	return nil
}
