package projects

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
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
	// Convert tagID to int
	id, err := strconv.Atoi(tagID)
	if err != nil {
		return nil, errors.New("invalid tag ID")
	}

	var tag Tag
	err = db.QueryRow(`SELECT id, name, color, projectid FROM TAG WHERE id = ?`, id).Scan(
		&tag.ID, &tag.Name, &tag.Color, &tag.ProjectID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tag not found")
		}
		return nil, err
	}

	// Convert to JSON
	jsonTag, err := json.Marshal(tag)
	if err != nil {
		return nil, err
	}

	return jsonTag, nil
}
