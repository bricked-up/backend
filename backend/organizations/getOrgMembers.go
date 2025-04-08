package organizations

import (
	"database/sql"
	"encoding/json"

	_ "modernc.org/sqlite"
)

// GetOrgMembers retrieves all member IDs belonging to a specific organization
func GetOrgMembers(db *sql.DB, orgID int) (string, error) {

	// Prepare query to get all members of the organization
	query := `
		SELECT userid 
		FROM ORG_MEMBER 
		WHERE orgid = ?
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
