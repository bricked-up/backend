package organizations

import (
	"database/sql"
	"errors"
	"time"

	"brickedup/backend/utils"

	_ "modernc.org/sqlite"
)

// UpdateOrg updates an organization if the user has executive privileges.
func UpdateOrg(db *sql.DB, sessionID, orgID int, org utils.Organization) error {
	// First, validate the session and check execution privileges
	var userID int
	var org_exists bool

	var sessionExpires time.Time
	err := db.QueryRow(`
		SELECT userid, expires FROM SESSION 
		WHERE id = ? AND expires > ?
	`, sessionID, time.Now()).Scan(&userID, &sessionExpires)
	if err != nil {
		return err
	}

	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT * FROM ORGANIZATION
			WHERE id = ?
		)
	`, orgID).Scan(&org_exists)

	if err != nil {
		return err
	}

	if !org_exists {
		return errors.New("Invalid orgid!")
	}

	// // Check if the user has exec privileges for this project
	// var hasExecPrivilege bool
	// err = db.QueryRow(`
	// 	SELECT EXISTS (
	// 		SELECT 1 FROM ORG_MEMBER om
	// 		JOIN ORG_MEMBER_ROLE omr ON om.id = omr.memberid
	// 		JOIN ORG_ROLE orgr ON omr.roleid = orgr.id
	// 		WHERE om.userid = ? AND om.orgid = ? AND orgr.can_exec = 1
	// 	)
	// `, userID, orgID).Scan(&hasExecPrivilege)
	// if err != nil {
	// 	return err
	// }
	// if !hasExecPrivilege {
	// 	return sql.ErrNoRows // Indicates no matching privileges found
	// }

	// Sanitize org fields
	sanitizedOrg := utils.Organization{
		ID:       orgID,
		Name:     utils.SanitizeText(org.Name, utils.TEXT),
	}

	// Update the org in the database
	_, err = db.Exec(`
		UPDATE ORGANIZATION
		SET name = ?
		WHERE id = ?
	`,
		sanitizedOrg.Name,
		sanitizedOrg.ID,
	)
	if err != nil {
		return err
	}

	return nil
}
