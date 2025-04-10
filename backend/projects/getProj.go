package projects

import (
	"brickedup/backend/utils"
	"database/sql"
	"encoding/json"
	"errors"

	_ "modernc.org/sqlite"
)

// GetProjMembers retrieves all member IDs belonging to a specific project.
func getProjMembers(db *sql.DB, projID int) ([]int, error) {
	query := `
		SELECT userid 
		FROM PROJECT_MEMBER 
		WHERE projectid = ?
	`

	// Execute the query
	rows, err := db.Query(query, projID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect all user IDs
	var memberIDs []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		memberIDs = append(memberIDs, userID)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return memberIDs, nil
}

// GetProjTags fetches all tags belonging to a project.
func getProjTags(db *sql.DB, projectid int) ([]int, error) {
	var tags []int

	rows, err := db.Query(
		`SELECT id
		FROM TAG 
		WHERE projectid = ?`, 
		projectid)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tagid int

		err = rows.Scan(&tagid)
		if err != nil {
			return nil, err
		}

		tags = append(tags, tagid)
	}

	return tags, nil
}

// GetProjIssues retrieves an array of all issues belonging to a project.
func getProjIssues(db *sql.DB, projectid int) ([]int, error) {
	var issues []int

	rows, err := db.Query(
		`SELECT issueid
		FROM PROJECT_ISSUES 
		WHERE projectid = ?`, 
		projectid)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var issueid int

		err = rows.Scan(&issueid)
		if err != nil {
			return nil, err
		}

		issues = append(issues, issueid)
	}

	return issues, nil
}

// GetProjRoles retrieves an array of all issues belonging to a project.
func getProjRoles(db *sql.DB, projectid int) ([]int, error) {
	var roles []int

	rows, err := db.Query(
		`SELECT id
		FROM PROJECT_ROLE
		WHERE projectid = ?`, 
		projectid)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var roleid int

		err = rows.Scan(&roleid)
		if err != nil {
			return nil, err
		}

		roles = append(roles, roleid)
	}

	return roles, nil
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

	project.Members, err = getProjMembers(db, projectID)
	if err != nil {
		return "", err
	}

	project.Tags, err = getProjTags(db, projectID)
	if err != nil {
		return "", err
	}

	project.Issues, err = getProjIssues(db, projectID)
	if err != nil {
		return "", err
	}

	project.Roles, err = getProjRoles(db, projectID)
	if err != nil {
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
