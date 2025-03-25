package backend

import (
	"database/sql"
	"errors"
)

// assignProjectRole promotes a validated user (userB) to a new role within a project,
// if the acting user (userA, identified via sessionid) has exec permission within that project.
func assignProjectRole(db *sql.DB, sessionid int, userB string, roleid, projectid int) error {
	// Validate session
	var userA int
	var sessionValid bool
	err := db.QueryRow(`
		SELECT userid, expires > CURRENT_TIMESTAMP 
		FROM SESSION WHERE id = ?
	`, sessionid).Scan(&userA, &sessionValid)
	if err != nil || !sessionValid {
		return errors.New("invalid session")
	}

	// Check project existence
	var exists bool
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM PROJECT WHERE id = ?)`, projectid).Scan(&exists)
	if err != nil || !exists {
		return errors.New("project not found")
	}

	// Check userB is verified
	var verified bool
	err = db.QueryRow(`SELECT verified FROM USER WHERE id = ?`, userB).Scan(&verified)
	if err != nil || !verified {
		return errors.New("userB is not verified")
	}

	// Get userB's project member id
	var userBMemberID int
	err = db.QueryRow(`
		SELECT id FROM PROJECT_MEMBER WHERE userid = ? AND projectid = ?
	`, userB, projectid).Scan(&userBMemberID)
	if err != nil {
		return errors.New("userB is not part of the project")
	}

	// Check userB does not already have the role
	err = db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM PROJECT_MEMBER_ROLE WHERE memberid = ? AND roleid = ?)
	`, userBMemberID, roleid).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("userB already has this role")
	}

	// Check if userA has exec permission
	var hasExec bool
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM PROJECT_MEMBER_ROLE pmr
			JOIN PROJECT_ROLE pr ON pmr.roleid = pr.id
			JOIN PROJECT_MEMBER pm ON pm.id = pmr.memberid
			WHERE pm.userid = ? AND pm.projectid = ? AND pr.can_exec = 1
		)
	`, userA, projectid).Scan(&hasExec)
	if err != nil || !hasExec {
		return errors.New("userA does not have exec permission")
	}

	// Insert new role
	_, err = db.Exec(`
		INSERT INTO PROJECT_MEMBER_ROLE (memberid, roleid) VALUES (?, ?)
	`, userBMemberID, roleid)
	if err != nil {
		return err
	}

	return nil
}
