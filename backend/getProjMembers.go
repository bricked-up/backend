package backend

import (
	"database/sql"
	"encoding/json"

	_ "modernc.org/sqlite"
)

// GetProjMembers retrieves all member IDs belonging to a specific project.
func GetProjMembers(db *sql.DB, orgID int) (string, error) {

	// Prepare query to get all members of the organization
	query := `
		SELECT userid 
		FROM PROJECT_MEMBER 
		WHERE projectid = ?
	`

	// Execute the query
	rows, err := db.Query(query, orgID)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	// Collect all user IDs
	var memberIDs []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return "", err
		}
		memberIDs = append(memberIDs, userID)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return "", err
	}

	// Marshal the result to JSON
	jsonResult, err := json.Marshal(memberIDs)
	if err != nil {
		return "", err
	}

	return string(jsonResult), nil
}
