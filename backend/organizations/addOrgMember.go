package organizations

import (
	"database/sql"
	"errors"
)

// AddOrgMember adds a user to a organization. The user is given a role within the organization.
func AddOrgMember(db *sql.DB, sessionid int64, userid int, roleid int, orgid int) error {
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
			SELECT 1 FROM ORG_MEMBER_ROLE omr
			JOIN ORG_ROLE orgr ON omr.roleid = orgr.id
			WHERE omr.memberid = (
				SELECT id FROM ORG_MEMBER WHERE userid = ? AND orgid = ?
			) AND orgr.can_exec = 1
		)`, manager, orgid).Scan(&has_exec)

	if err != nil {
		return err
	}

	if !has_exec {
		return errors.New("User does not have exec privileges in the organization!")
	}

	err = db.QueryRow(`
		SELECT EXISTS (
		    SELECT 1
    		FROM USER
    		WHERE id = ? 
		)
	`, userid).Scan(&user_exists)

	if !user_exists {
		return errors.New("User does not exist! Cannot add to organization!")
	}

	res, err := db.Exec(`
		INSERT INTO PROJECT_MEMBER(userid, projectid)
		VALUES(?, ?)
	`, userid, orgid)

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
