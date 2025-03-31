package backend

import (
	"database/sql"
	"errors"

	"brickedup/backend/utils"
)

func UpdateOrganizationName(db *sql.DB, sessionID int, orgID int, newName string) error {
	if orgID <= 0 {
		return errors.New("invalid organization ID")
	}

	// Sanitize the organization name
	sanitizedName := utils.SanitizeText(newName, utils.TEXT)
	if sanitizedName == "" {
		return errors.New("organization name cannot be empty")
	}

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Check if the organization exists
	var exists bool
	err = tx.QueryRow("SELECT EXISTS (SELECT 1 FROM ORGANIZATION WHERE id = ?)", orgID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("organization not found")
	}

	// Check if the user has exec privileges
	var hasExecPriv bool
	err = tx.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM SESSION s
			JOIN ORG_MEMBER om ON s.userid = om.userid
			JOIN ORG_MEMBER_ROLE omr ON om.id = omr.memberid
			JOIN ORG_ROLE r ON omr.roleid = r.id
			WHERE s.id = ? AND om.orgid = ? AND r.can_exec = 1
		)
	`, sessionID, orgID).Scan(&hasExecPriv)
	if err != nil {
		return err
	}
	if !hasExecPriv {
		return errors.New("not authorized")
	}

	// Check for name conflicts
	var nameExists bool
	err = tx.QueryRow("SELECT EXISTS (SELECT 1 FROM ORGANIZATION WHERE name = ? AND id != ?)", sanitizedName, orgID).Scan(&nameExists)
	if err != nil {
		return err
	}
	if nameExists {
		return errors.New("organization name already exists")
	}

	// Update the organization name
	result, err := tx.Exec("UPDATE ORGANIZATION SET name = ? WHERE id = ?", sanitizedName, orgID)
	if err != nil {
		return err
	}

	// Ensure one row was updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return errors.New("unexpected number of rows updated")
	}

	return nil
}
