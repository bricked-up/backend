package issues

import (
	"database/sql"
	"encoding/json"

	_ "modernc.org/sqlite"
)

// GetIssue fetches issue details and returns them as a JSON string
func GetIssue(db *sql.DB, issueid int) (string, error) {

	row := db.QueryRow("SELECT id, title, desc, tagid, priority, created, completed, cost FROM ISSUE WHERE id = ?", issueid)

	// Variables to store fetched data
	var id, tagid, priority, cost int
	var title, description string
	var created, completed sql.NullTime

	// Scan row into variables
	err := row.Scan(&id, &title, &description, &tagid, &priority, &created, &completed, &cost)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", err
		}
		return "", err
	}

	// Convert the issue details into a map
	issueDetails := map[string]interface{}{
		"id":          id,
		"title":       title,
		"description": description,
		"tagid":       tagid,
		"priority":    priority,
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
