package organizations

import (
	"brickedup/backend/utils"
	"database/sql"
	"encoding/json"
	"errors"

	_ "modernc.org/sqlite"
)

// GetOrg returns a JSON formatted string of an organization entry.
func GetOrg(db *sql.DB, orgid int) ([]byte, error) {

	row := db.QueryRow(`SELECT id, name FROM organization where id = ?`, orgid)

	var org utils.Organization
	if err := row.Scan(&org.ID, &org.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Organization not found")
		}
		return nil, err
	}

	jsonOrg, err := json.Marshal(org)
	if err != nil {
		return nil, err
	}

	return jsonOrg, err
}
