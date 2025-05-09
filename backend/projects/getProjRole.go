package projects

import (
	"brickedup/backend/utils"
	"database/sql"
)

// GetProjRole gets returns the ProjectRole corresponding to roleid.
func GetProjRole(db *sql.DB, roleid int) (*utils.ProjectRole, error) {
	role := &utils.ProjectRole{}
	role.ID = roleid

	err := db.QueryRow(`
		SELECT projectid, name, can_exec, can_write, can_read
		FROM PROJECT_ROLE
		WHERE id = ?
	`, role.ID).Scan(
		&role.ProjectID, 
		&role.Name,
		&role.CanExec,
		&role.CanWrite,
		&role.CanRead)

	if err != nil {
		return nil, err
	}

	return role, nil
}
