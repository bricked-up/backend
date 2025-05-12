package projects

import (
	"database/sql"
	"time"

	"brickedup/backend/utils"

	_ "modernc.org/sqlite"
)

// UpdateProject function updates project.
func UpdateProject(db *sql.DB, sessionID, projectID int, project utils.Project) error {
	// First, validate the session and check execution privileges
	var userID int
	var sessionExpires time.Time
	err := db.QueryRow(`
		SELECT userid, expires FROM SESSION 
		WHERE id = ? AND expires > ?
	`, sessionID, time.Now()).Scan(&userID, &sessionExpires)
	if err != nil {
		return err
	}

	// Check if the user has exec privileges for this project
	var hasExecPrivilege bool
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM PROJECT_MEMBER pm
			JOIN PROJECT_MEMBER_ROLE pmr ON pm.id = pmr.memberid
			JOIN PROJECT_ROLE pr ON pmr.roleid = pr.id
			WHERE pm.userid = ? AND pm.projectid = ? AND pr.can_exec = 1
		)
	`, userID, projectID).Scan(&hasExecPrivilege)
	if err != nil {
		return err
	}
	if !hasExecPrivilege {
		return sql.ErrNoRows // Indicates no matching privileges found
	}

	// Sanitize project fields
	sanitizedProject := utils.Project{
		ID:       project.ID,
		OrgID:    project.OrgID,
		Name:     utils.SanitizeText(project.Name, utils.TEXT),
		Budget:   project.Budget,
		Charter:  utils.SanitizeText(project.Charter, utils.TEXT),
		Archived: project.Archived,
	}

	// Update the project in the database
	result, err := db.Exec(`
		UPDATE PROJECT 
		SET name = ?, budget = ?, charter = ?, archived = ?
		WHERE id = ? AND orgid = ?
	`,
		sanitizedProject.Name,
		sanitizedProject.Budget,
		sanitizedProject.Charter,
		sanitizedProject.Archived,
		projectID,
		sanitizedProject.OrgID,
	)
	if err != nil {
		return err
	}

	// Check if any rows were actually updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
