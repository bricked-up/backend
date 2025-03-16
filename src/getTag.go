package backend

import (
	"database/sql"
	"encoding/json"
	"errors"
)

// Tag holds the details for a tag.
type Tag struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	ProjectID int    `json:"project_id"`
}

// getTagDetails fetches tag details by tag ID and returns JSON data.
func getTagDetails(db *sql.DB, tagID string) ([]byte, error) {
	var tag Tag

	err := db.QueryRow(`SELECT id, name, color, projectid FROM TAG WHERE id = ?`, tagID).Scan(
		&tag.ID, &tag.Name, &tag.Color, &tag.ProjectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Tag not found")
		}
		return nil, err
	}

	// Convert tag data to JSON
	jsonTag, err := json.Marshal(tag)
	if err != nil {
		return nil, err
	}

	return jsonTag, nil
}
