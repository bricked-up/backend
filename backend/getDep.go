package backend

import (
	"database/sql"
	"encoding/json"
	"errors"

	_ "modernc.org/sqlite"
)

type Dep struct {
	ID int `json:"id"`
}

// GetDep getches the relations of dependencies of one issue and returns it as a json data
func getDep(db *sql.DB, issueid int) ([]byte, error) {
	row := db.QueryRow(`SELECT dependency FROM DEPENDENCY WHERE issueid = ?`, issueid)
	var dep Dep
	// Scans the Rows to see if there are any
	if err := row.Scan(&dep.ID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("issueid not found")
		}
	}
	//Convert the Dep struct to JSON
	jsonDep, err := json.Marshal(dep)
	if err != nil {
		return nil, err
	}

	return jsonDep, nil
}
