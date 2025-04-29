package projects

import (
	"database/sql"
	"errors"
)

// ArchiveProj marks a project as "archived" if the user has the necessary
// privileges in the project's organization.
func ArchiveProj(db *sql.DB, sessionid int, projectid int) error {
	var userid int
	var orgid int
	var has_exec bool

	err := db.QueryRow(
		`SELECT userid FROM SESSION WHERE id = ?`,
		sessionid).Scan(&userid)

	if err != nil {
		return err
	}

	err = db.QueryRow(
		`SELECT orgid FROM PROJECT WHERE id = ?`,
		projectid).Scan(&orgid)

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
		)`, userid, orgid).Scan(&has_exec)

	if err != nil || !has_exec {
		return errors.New("user lacks exec permissions")
	}

	_, err = db.Exec(
		`UPDATE PROJECT SET archived = 1 WHERE id = ?`, 
		projectid)

	if err != nil {
		return err
	}

	return nil
}
