package organizations

import (
	"brickedup/backend/utils"
	"database/sql"
)

// GetOrgRole returns the OrgRole corresponding to roleid.
func GetOrgRole(db *sql.DB, roleid int) (*utils.OrgRole, error) {
	role := &utils.OrgRole{}
	role.ID = roleid

	err := db.QueryRow(`
		SELECT orgid, name, can_exec, can_write, can_read
		FROM ORG_ROLE
		WHERE id = ?
	`, role.ID).Scan(
		&role.OrgID, 
		&role.Name,
		&role.CanExec,
		&role.CanWrite,
		&role.CanRead)

	if err != nil {
		return nil, err
	}

	return role, nil
}
