package organizations

import (
	"database/sql"
	"errors"
)

// RemoveOrgMember removes a user from an organization.
// The user is given a role within the organization.
func RemoveOrgMember(db *sql.DB, sessionid int64, memberid int) error {
	var manager int
	var has_exec bool

	err := db.QueryRow(
		`SELECT userid FROM SESSION WHERE id = ?`,
		sessionid).Scan(&manager)

	if err != nil {
		return err
	}

	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT * FROM ORG_MEMBER_ROLE omr
			JOIN ORG_ROLE orgr ON omr.roleid = orgr.id
			WHERE omr.memberid = (
				SELECT id
				FROM ORG_MEMBER
				WHERE userid = (
					SELECT userid
					FROM SESSION
					WHERE id = ?
				)
			) AND orgr.can_exec = 1
		)`, sessionid).Scan(&has_exec)

	if err != nil {
		return err
	}

	if !has_exec {
		return errors.New("User does not have exec privileges in the organization!")
	}

	_, err = db.Exec(`
		DELETE FROM ORG_MEMBER
		WHERE id = ?
	`, memberid)

	if err != nil {
		return err
	}

	_, err = db.Exec(`
		DELETE FROM ORG_MEMBER_ROLE
		WHERE memberid = ?
	`, memberid)

	if err != nil {
		return err
	}

	return nil
}
