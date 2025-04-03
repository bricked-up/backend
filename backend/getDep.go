package backend

import (
	"database/sql"
	"encoding/json"

	_ "modernc.org/sqlite"
)

// GetDep getches the relations of dependencies of one issue and returns it as a json data
func getDep(db *sql.DB, issueid int) ([]byte, error) {
	rows, err := db.Query(`SELECT dependency FROM DEPENDENCY WHERE issueid = ?`, issueid)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var dependencies []int64
	for rows.Next() {
		var dep int64

		err := rows.Scan(&dep)
		if err != nil {
			return nil, err
		}

		dependencies = append(dependencies, dep)
	}

	//Convert the Dep struct to JSON
	jsonDep, err := json.Marshal(dependencies)
	if err != nil {
		return nil, err
	}

	return jsonDep, nil
}
