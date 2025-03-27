package backend

import (
	"database/sql"
	"encoding/json"

	_ "modernc.org/sqlite"
)

// GetIssueDetails fetches issue details and returns them as a JSON string
func getIssueDetails(db *sql.DB, issueid int) (string, error) {

<<<<<<< HEAD
	row := db.QueryRow("SELECT id, title, desc, tagid, priority, created, completed, cost FROM ISSUE WHERE id = ?", issueid)

	// Variables to store fetched data
	var id, tagid, priority, cost int
=======
	row := db.QueryRow("SELECT id, title, desc, tagid, priorityid, created, completed, cost FROM ISSUE WHERE id = ?", issueid)

	// Variables to store fetched data
	var id, tagid, priorityid, cost int
>>>>>>> 462d0b2 (bood)
	var title, description string
	var created, completed sql.NullTime

	// Scan row into variables
<<<<<<< HEAD
	err := row.Scan(&id, &title, &description, &tagid, &priority, &created, &completed, &cost)
=======
	err := row.Scan(&id, &title, &description, &tagid, &priorityid, &created, &completed, &cost)
>>>>>>> 462d0b2 (bood)
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
<<<<<<< HEAD
		"priority":  priority,
=======
		"priorityid":  priorityid,
>>>>>>> 462d0b2 (bood)
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
