package issues

import (
	"brickedup/backend/utils"
	"database/sql"
	"encoding/json"

	_ "modernc.org/sqlite"
)

// GetIssueDep fetches all issues that the issue depends on.
func getIssueDep(db *sql.DB, issue *utils.Issue) error {
	rows, err := db.Query(
		`SELECT dependency FROM DEPENDENCY WHERE issueid = ?`,
		issue.ID)

	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var dep int

		err := rows.Scan(&dep)
		if err != nil {
			return err
		}

		issue.Dependencies = append(issue.Dependencies, dep)
	}

	return nil
}

// GetIssue fetches issue details and returns them as a JSON string
func GetIssue(db *sql.DB, issueid int) (string, error) {
	row := db.QueryRow("SELECT title, desc, tagid, priority, created, completed, cost FROM ISSUE WHERE id = ?", issueid)

	var issue utils.Issue
	issue.ID = issueid

	// Scan row into variables
	err := row.Scan(
		&issue.Title, 
		&issue.Desc, 
		&issue.TagID, 
		&issue.Priority, 
		&issue.Created, 
		&issue.Completed, 
		&issue.Cost)

	if err != nil {
		return "", err
	}

	// Convert map to JSON
	jsonData, err := json.Marshal(issue)
	if err != nil {
		return "", err
	}

	// Return JSON as string
	return string(jsonData), nil
}
