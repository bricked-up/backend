package projects

import (
	"brickedup/backend/utils"
	"database/sql"
	"errors"
)

// GetTag fetches tag details by tag ID and returns JSON data.
func GetTag(db *sql.DB, tagID int) (*utils.Tag, error) {
	tag := &utils.Tag{}
	tag.ID = tagID

	err := db.QueryRow(`SELECT id, name, color, projectid FROM TAG WHERE id = ?`, tag.ID).Scan(
		&tag.ID, &tag.Name, &tag.Color, &tag.ProjectID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tag not found")
		}
		return nil, err
	}

	return tag, nil
}
