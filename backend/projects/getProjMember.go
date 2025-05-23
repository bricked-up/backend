package projects

import (
	"brickedup/backend/utils"
	"database/sql"

	_ "modernc.org/sqlite"
)

// Gets all issues that have been assigned to the member in the project.
func getMemberIssues(db *sql.DB, member *utils.ProjectMember) error {
	member.Issues = nil

	rows, err := db.Query(
		`SELECT DISTINCT ui.issueid
		FROM USER_ISSUES ui
		JOIN PROJECT_ISSUES pi ON ui.issueid = pi.issueid
		WHERE ui.userid = ?`, member.ID)
	if err != nil {
		return err
	}

	for rows.Next() {
		var issueid int

		err := rows.Scan(&issueid)
		if err != nil {
			return err
		}

		member.Issues = append(member.Issues, issueid)
	}

	return nil
}

// GetMemberRoles retrieves the roles of the member as well as sets user privileges.
func getMemberRoles(db *sql.DB, member *utils.ProjectMember) error {
	rows, err := db.Query(
		`SELECT roleid
		FROM PROJECT_MEMBER_ROLE
		WHERE memberid = ?`, member.ID)

	if err != nil {
		return err
	}

	for rows.Next() {
		var roleid int

		err := rows.Scan(&roleid)
		if err != nil {
			return err
		}

		member.Roles = append(member.Roles, roleid)
	}

	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT *
			FROM PROJECT_MEMBER_ROLE pmr
			JOIN PROJECT_ROLE pr ON pmr.roleid = pr.id
			WHERE pmr.memberid = ? AND pr.can_exec = 1 
		)`, member.ID).Scan(&member.CanExec)

	if err != nil {
		return err
	}

	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT *
			FROM PROJECT_MEMBER_ROLE pmr
			JOIN PROJECT_ROLE pr ON pmr.roleid = pr.id
			WHERE pmr.memberid = ? AND pr.can_write = 1 
		)`, member.ID).Scan(&member.CanWrite)

	if err != nil {
		return err
	}

	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT *
			FROM PROJECT_MEMBER_ROLE pmr
			JOIN PROJECT_ROLE pr ON pmr.roleid = pr.id
			WHERE pmr.memberid = ? AND pr.can_read = 1 
		)`, member.ID).Scan(&member.CanRead)

	if err != nil {
		return err
	}

	return nil
}


// GetProjMember fetches a project member by its memberid from the DB and 
// returns a ProjectMember struct.
func GetProjMember(db *sql.DB, memberid int) (*utils.ProjectMember, error) {
	member := &utils.ProjectMember{}
	member.ID = memberid

	row := db.QueryRow(
		`SELECT userid, projectid 
		FROM PROJECT_MEMBER
		WHERE id = ?`, member.ID)

	err := row.Scan(&member.UserID, &member.ProjectID)
	if err != nil {
		return nil, err
	}

	err = getMemberIssues(db, member)
	if err != nil {
		return nil, err
	}

	err = getMemberRoles(db, member)
	if err != nil {
		return nil, err
	}

	return member, nil
}
