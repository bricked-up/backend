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

func getDep(db *sql.DB, issueid int) ([]byte, error) {
	row := db.QueryRow(`SELECT dependency FROM DEPENDENCY WHERE issueid = ?`, issueid)
	var dep Dep
	if err := row.Scan(&dep.ID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("issueid not found")
		}
	}

	jsonDep, err := json.Marshal(dep)
	if err != nil {
		return nil, err
	}

	return jsonDep, nil
}
