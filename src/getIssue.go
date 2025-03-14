package backend

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "modernc.org/sqlite"
)

// GetIssueDetails fetches issue details and returns them as a JSON string
func getIssueDetails(db *sql.DB, issueid int) (string, error) {

	row := db.QueryRow("SELECT id, title, desc, tagid, priorityid, created, completed, cost FROM ISSUE WHERE id = ?", issueid)

	// Variables to store fetched data
	var id, tagid, priorityid, cost int
	var title, description string
	var created, completed sql.NullTime

	// Scan row into variables
	err := row.Scan(&id, &title, &description, &tagid, &priorityid, &created, &completed, &cost)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no issue found with ID %d", issueid)
		}
		return "", err
	}

	// Convert the issue details into a map
	issueDetails := map[string]interface{}{
		"id":          id,
		"title":       title,
		"description": description,
		"tagid":       tagid,
		"priorityid":  priorityid,
		"created":     created.Time,
		"completed":   completed.Time,
		"cost":        cost,
	}

	// Convert map to JSON
	jsonData, err := json.Marshal(issueDetails)
	if err != nil {
		return "", err
	}

	// Return JSON as string
	return string(jsonData), nil
}
