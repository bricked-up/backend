package backend

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

// RemoveOrgMemberRole removes a role from a user within an organization.
func RemoveOrgMemberRole(db *sql.DB, sessionid int, orgMemberRoleId int) error {
	// Validate user A's session and get their user ID
	var userAID int
	err := db.QueryRow("SELECT userid FROM SESSION WHERE id = ?", sessionid).Scan(&userAID)
	if err != nil {
		return err
	}

	// Get the organization ID for the role being removed
	var orgID int
	err = db.QueryRow(`
        SELECT o.id 
        FROM ORG_MEMBER_ROLE omr
        JOIN ORG_ROLE r ON omr.roleid = r.id
        JOIN ORG_MEMBER m ON omr.memberid = m.id
        JOIN ORGANIZATION o ON r.orgid = o.id
        WHERE omr.id = ?
    `, orgMemberRoleId).Scan(&orgID)
	if err != nil {
		return err
	}

	// Check if user A has Admin role in the organization
	var hasPermission bool
	err = db.QueryRow(`
        SELECT EXISTS (
            SELECT 1 
            FROM ORG_MEMBER_ROLE omr
            JOIN ORG_ROLE r ON omr.roleid = r.id
            JOIN ORG_MEMBER m ON omr.memberid = m.id
            WHERE m.userid = ? AND r.name = 'Admin' AND r.orgid = ?
        )
    `, userAID, orgID).Scan(&hasPermission)
	if err != nil {
		return err
	}
	if !hasPermission {
		return err
	}

	// Remove the specified role assignment from ORG_MEMBER_ROLE
	result, err := db.Exec("DELETE FROM ORG_MEMBER_ROLE WHERE id = ?", orgMemberRoleId)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return err
	}
	return nil
}
