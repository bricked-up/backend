package projects

import (
	"database/sql"
	"errors"
)

// AddProjMember adds a user to a project. The user is given a role within the project.
func AddProjMember(db *sql.DB, sessionid int64, userid int, roleid int, projectid int) error {
	var manager int
	var has_exec bool
	var user_exists bool

	err := db.QueryRow(
		`SELECT userid FROM SESSION WHERE id = ?`,
		sessionid).Scan(&manager)

	if err != nil {
		return err
	}

	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM PROJECT_MEMBER_ROLE pmr
			JOIN PROJECT_ROLE pr ON pmr.roleid = pr.id
			WHERE pmr.memberid = (
				SELECT id FROM PROJECT_MEMBER WHERE userid = ? AND projectid = ?
			) AND pr.can_exec = 1
		)`, manager, projectid).Scan(&has_exec)

	if err != nil {
		return err
	}

	if !has_exec {
		return errors.New("User does not have exec privileges in the project!")
	}

	err = db.QueryRow(`
		SELECT EXISTS (
		    SELECT 1
    		FROM USER
    		WHERE id = ? 
		)
	`, userid).Scan(&user_exists)

	if !user_exists {
		return errors.New("User does not exist! Cannot add to project!")
	}

	res, err := db.Exec(`
		INSERT INTO PROJECT_MEMBER(userid, projectid)
		VALUES(?, ?)
	`, userid, projectid)

	memberid, err := res.LastInsertId()
	if err != nil {
		return err
	}

	res, err = db.Exec(`
		INSERT INTO PROJECT_MEMBER_ROLE(memberid, roleid)
		VALUES(?, ?)
	`, memberid, roleid)

	if err != nil {
		return err
	}

	return nil
}
