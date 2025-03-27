package backend

import (
	"database/sql"
	"encoding/json"
	"errors"

	_ "modernc.org/sqlite"
)

type Org struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetOrg(db *sql.DB, orgid int) ([]byte, error) {

	row := db.QueryRow(`SELECT id, name FROM organization where id = ?`, orgid)

	var org Org
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
