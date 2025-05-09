package organizations

import (
	"brickedup/backend/utils"
	"database/sql"

	_ "modernc.org/sqlite"
)

// GetMemberRoles retrieves the roles of the member as well as sets user privileges.
func getMemberRoles(db *sql.DB, member *utils.OrgMember) error {
	rows, err := db.Query(
		`SELECT roleid
		FROM ORG_MEMBER_ROLE
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

	row := db.QueryRow(`
		SELECT 
			MAX(CASE WHEN orgr.can_read THEN 1 ELSE 0 END),
			MAX(CASE WHEN orgr.can_write THEN 1 ELSE 0 END),
			MAX(CASE WHEN orgr.can_exec THEN 1 ELSE 0 END)
		FROM ORG_MEMBER_ROLE omr
		JOIN ORG_ROLE orgr ON omr.roleid = orgr.id
		WHERE omr.memberid = ? `, member.ID)

	err = row.Scan(&member.CanRead, &member.CanWrite, &member.CanExec)
	if err != nil {
		return err
	}

	return nil
}


// GetOrgMember fetches an organization member by its memberid from the DB and 
// returns an OrgMember struct.
func GetOrgMember(db *sql.DB, memberid int) (*utils.OrgMember, error) {
	member := &utils.OrgMember{}
	member.ID = memberid

	row := db.QueryRow(
		`SELECT userid, orgid 
		FROM ORG_MEMBER
		WHERE id = ?`, member.ID)

	err := row.Scan(&member.UserID, &member.OrganizationID)
	if err != nil {
		return nil, err
	}

	err = getMemberRoles(db, member)
	if err != nil {
		return nil, err
	}

	return member, nil
}
