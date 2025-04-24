package projects

import (
	"brickedup/backend/utils"
	"database/sql"
	"encoding/json"
	"errors"

	_ "modernc.org/sqlite"
)

// GetProjMembers retrieves all member IDs belonging to a specific project.
func getProjMembers(db *sql.DB, proj *utils.Project) error {
	proj.Members = nil

	query := `
		SELECT userid 
		FROM PROJECT_MEMBER 
		WHERE projectid = ?
	`

	// Execute the query
	rows, err := db.Query(query, proj.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Collect all user IDs
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return err
		}
		proj.Members = append(proj.Members, userID)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

// GetProjTags fetches all tags belonging to a project.
func getProjTags(db *sql.DB, proj *utils.Project) error {
	proj.Tags = nil

	rows, err := db.Query(
		`SELECT id
		FROM TAG 
		WHERE projectid = ?`, 
		proj.ID)

	if err != nil {
		return err
	}

	for rows.Next() {
		var tagid int

		err = rows.Scan(&tagid)
		if err != nil {
			return err
		}

		proj.Tags = append(proj.Tags, tagid)
	}

	return nil
}

// GetProjIssues retrieves an array of all issues belonging to a project.
func getProjIssues(db *sql.DB, proj *utils.Project) error {
	proj.Issues = nil

	rows, err := db.Query(
		`SELECT issueid
		FROM PROJECT_ISSUES 
		WHERE projectid = ?`, 
		proj.ID)

	if err != nil {
		return err
	}

	for rows.Next() {
		var issueid int

		err = rows.Scan(&issueid)
		if err != nil {
			return err
		}

		proj.Issues = append(proj.Issues, issueid)
	}

	return nil
}

// GetProjRoles retrieves an array of all issues belonging to a project.
func getProjRoles(db *sql.DB, proj *utils.Project) error {
	proj.Roles = nil

	rows, err := db.Query(
		`SELECT id
		FROM PROJECT_ROLE
		WHERE projectid = ?`, 
		proj.ID)

	if err != nil {
		return err
	}

	for rows.Next() {
		var roleid int

		err = rows.Scan(&roleid)
		if err != nil {
			return err
		}

		proj.Roles = append(proj.Roles, roleid)
	}

	return nil
}

// GetProject fetches all the details of a project and returns them as a JSON string
func GetProject(db *sql.DB, projectID int) (string, error) {
	// validate projectID is not null or negative value
	if projectID <= 0 {
		return "", errors.New("invalid project ID")
	}

	// Perform query using the projectID (no need to sanitize for an integer)
	row := db.QueryRow(
		`SELECT orgid, name, budget, charter, archived 
		FROM PROJECT 
		WHERE id = ?`, 
		projectID)

	var project utils.Project
	project.ID = projectID

	// Scan row into variables
	err := row.Scan(
		&project.OrgID, 
		&project.Name, 
		&project.Budget, 
		&project.Charter, 
		&project.Archived)

	if err != nil {
		return "", err
	}

	if err := getProjMembers(db, &project); err != nil {
		return "", err
	}

	if err := getProjTags(db, &project); err != nil {
		return "", err
	}

	if err := getProjIssues(db, &project); err != nil {
		return "", err
	}

	if err := getProjRoles(db, &project); err != nil {
		return "", err
	}

	// Convert map to JSON
	jsonData, err := json.Marshal(project)
	if err != nil {
		return "", err
	}

	// Return JSON as string
	return string(jsonData), nil
}
