package backend

import (
	"context"
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

// Project represents a project in the database
type Project struct {
	ID       int    `json:"id"`
	OrgID    int    `json:"orgid"`
	Name     string `json:"name"`
	Budget   int    `json:"budget"`
	Charter  string `json:"charter"`
	Archived bool   `json:"archived"`
}

// UpdateProject updates all non-foreign key fields of a project
func updateProject(ctx context.Context, db *sql.DB, sessionID int, projectID int, project Project) error {
	// Validate that the project exists
	var existingProjectID int
	err := db.QueryRowContext(ctx, "SELECT id FROM PROJECT WHERE id = ?", projectID).Scan(&existingProjectID)
	if err != nil {
		return err
	}

	// Get the user ID from the session
	var userID int
	err = db.QueryRowContext(ctx, "SELECT userid FROM SESSION WHERE id = ? AND expires > ?",
		sessionID, time.Now()).Scan(&userID)
	if err != nil {
		return err
	}

	// Check if the user has exec privileges for the project
	var hasExecPrivilege bool
	err = db.QueryRowContext(ctx, `
        SELECT EXISTS (
            SELECT 1 FROM PROJECT_MEMBER pm
            JOIN PROJECT_MEMBER_ROLE pmr ON pm.id = pmr.memberid
            JOIN PROJECT_ROLE pr ON pmr.roleid = pr.id
            WHERE pm.userid = ? AND pm.projectid = ? AND pr.can_exec = 1
        )`, userID, projectID).Scan(&hasExecPrivilege)
	if err != nil {
		return err
	}

	if !hasExecPrivilege {
		return err
	}

	// Update the non-foreign key fields
	_, err = db.ExecContext(ctx, `
        UPDATE PROJECT
        SET name = ?, budget = ?, charter = ?, archived = ?
        WHERE id = ?`,
		project.Name, project.Budget, project.Charter, project.Archived, projectID)
	if err != nil {
		return err
	}

	return nil
}
