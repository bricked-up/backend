package projects

import (
	"database/sql"
	"encoding/json"
	"errors"

	_ "modernc.org/sqlite"
)

// GetProjectDetails fetches all the details of a project and returns them as a JSON string
func GetProjectDetails(db *sql.DB, projectID int) (string, error) {

	// validate projectID is not null or negative value
	if projectID <= 0 {
		return "", errors.New("invalid project ID")
	}

	// Perform query using the projectID (no need to sanitize for an integer)
	row := db.QueryRow("SELECT id, orgid, name, budget, charter, archived FROM PROJECT WHERE id = ?", projectID)

	// Variables to store fetched data
	var id, orgid, budget int
	var name, charter string
	var archived bool

	// Scan row into variables
	err := row.Scan(&id, &orgid, &name, &budget, &charter, &archived)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("project not found")
		}
		return "", err
	}

	// Convert the project details into a map
	projectDetails := map[string]interface{}{
		"id":       id,
		"orgid":    orgid,
		"name":     name,
		"budget":   budget,
		"charter":  charter,
		"archived": archived,
	}

	// Convert map to JSON
	jsonData, err := json.Marshal(projectDetails)
	if err != nil {
		return "", err
	}

	// Return JSON as string
	return string(jsonData), nil
}
