package backend

import (
	"database/sql"
	"errors"
)

func promoteUserRole(db *sql.DB, projectID, userID, roleID, userRoleID int) error {
	var userExists bool
	var userValidated bool
	var userInProject bool
	var userRoleExec bool

	// Check if user exists
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM USER WHERE id = ?)`, userID).Scan(&userExists)
	if err != nil || !userExists {
		return errors.New("user B does not exist")
	}

	// Check if user is validated
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM VERIFY_USER WHERE id = (SELECT verifyid FROM USER WHERE id = ?))`, userID).Scan(&userValidated)
	if err != nil || !userValidated {
		return errors.New("user B is not validated")
	}

	// Check if user is part of the project
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM PROJECT_MEMBER WHERE userid = ? AND projectid = ?)`, userID, projectID).Scan(&userInProject)
	if err != nil || !userInProject {
		return errors.New("user B is not part of the project")
	}

	// Check if user has the necessary exec permissions
	err = db.QueryRow(`SELECT can_exec FROM PROJECT_MEMBER_ROLE WHERE memberid = (SELECT id FROM PROJECT_MEMBER WHERE userid = ? AND projectid = ?) AND roleid = ?`, userID, projectID, userRoleID).Scan(&userRoleExec)
	if err != nil || !userRoleExec {
		return errors.New("user A lacks exec permissions")
	}

	// Check if user already has the specified role
	var currentRoleID int
	err = db.QueryRow(`SELECT roleid FROM PROJECT_MEMBER_ROLE WHERE memberid = (SELECT id FROM PROJECT_MEMBER WHERE userid = ? AND projectid = ?)`, userID, projectID).Scan(&currentRoleID)
	if err != nil {
		return err
	}
	if currentRoleID == roleID {
		return errors.New("user B already has the specified role")
	}

	// Promote user to the new role
	_, err = db.Exec(`UPDATE PROJECT_MEMBER_ROLE SET roleid = ? WHERE memberid = (SELECT id FROM PROJECT_MEMBER WHERE userid = ? AND projectid = ?)`, roleID, userID, projectID)
	return err
}
