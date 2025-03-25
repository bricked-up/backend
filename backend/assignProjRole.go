package backend

import (
	"database/sql"
	"errors"
)

// PromoteUserRole promotes a user (userID) to a new role (roleID) within a project,
// only if another user (userAID) with exec permission initiates the change.
func promoteUserRole(db *sql.DB, userAID, userBID, roleID, projectID int) error {
	// Check if user B exists
	var exists bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM USER WHERE id = ?)`, userBID).Scan(&exists)
	if err != nil || !exists {
		return errors.New("user B does not exist")
	}

	// Check if user B is validated
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM USER
			WHERE id = ? AND verifyid IS NOT NULL
		)
	`, userBID).Scan(&exists)
	if err != nil || !exists {
		return errors.New("user B is not validated")
	}

	// Check if user B is a member of the project
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM PROJECT_MEMBER
			WHERE userid = ? AND projectid = ?
		)
	`, userBID, projectID).Scan(&exists)
	if err != nil || !exists {
		return errors.New("user B is not part of the project")
	}

	// Check if user A has exec permissions
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM PROJECT_MEMBER_ROLE pmr
			JOIN PROJECT_MEMBER pm ON pm.id = pmr.memberid
			JOIN PROJECT_ROLE pr ON pr.id = pmr.roleid
			WHERE pm.userid = ? AND pm.projectid = ? AND pr.can_exec = 1
		)
	`, userAID, projectID).Scan(&exists)
	if err != nil || !exists {
		return errors.New("user A lacks exec permissions")
	}

	// Get user B's member ID in the project
	var memberID int
	err = db.QueryRow(`
		SELECT id FROM PROJECT_MEMBER
		WHERE userid = ? AND projectid = ?
	`, userBID, projectID).Scan(&memberID)
	if err != nil {
		return err
	}

	// Check if user B already has the specified role
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM PROJECT_MEMBER_ROLE
			WHERE memberid = ? AND roleid = ?
		)
	`, memberID, roleID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user B already has the specified role")
	}

	// Perform role update
	_, err = db.Exec(`
		UPDATE PROJECT_MEMBER_ROLE
		SET roleid = ?
		WHERE memberid = ?
	`, roleID, memberID)
	if err != nil {
		return err
	}

	return nil
}
