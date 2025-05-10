package projects

import (
	"database/sql"
	"errors"
)

func RemoveProjMember(db *sql.DB, sessionid int64, memberid int) error {
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
			SELECT 1 FROM PROJECT_MEMBER_ROLE pmr
			JOIN PROJECT_ROLE pr ON pmr.roleid = pr.id
			WHERE pmr.memberid = ? AND pr.can_exec = 1
		)`, manager, memberid).Scan(&has_exec)

	if err != nil {
		return err
	}

	if !has_exec {
		return errors.New("User does not have exec privileges in the project!")
	}

	_, err = db.Exec(`
		DELETE FROM PROJECT_MEMBER
		WHERE id = ?
	`, memberid)

	_, err = db.Exec(`
		DELETE FROM PROJECT_MEMBER_ROLE
		WHERE memberid = ?
	`, memberid)

	if err != nil {
		return err
	}

	return nil
}

