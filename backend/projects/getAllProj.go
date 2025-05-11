package projects

import "database/sql"

// GetAllProj retrieves the IDs of all the projects.
func GetAllProj(db *sql.DB) ([]int, error) {
	var projectids []int

	res, err := db.Query(`
		SELECT id FROM PROJECT
		`)

	if err != nil {
		return nil, err
	}

	for res.Next() {
		var projid int
		err = res.Scan(&projid)

		if err != nil {
			return nil, err
		}

		projectids = append(projectids, projid)
	}

	return projectids, nil
}
