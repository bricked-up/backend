package organizations

import "database/sql"

// GetAllOrg retrieves the IDs of all the organizations.
func GetAllOrg(db *sql.DB) ([]int, error) {
	var orgids []int

	res, err := db.Query(`
		SELECT id FROM ORGANIZATION
		`)

	if err != nil {
		return nil, err
	}

	for res.Next() {
		var orgid int
		err = res.Scan(&orgid)

		if err != nil {
			return nil, err
		}

		orgids = append(orgids, orgid)
	}

	return orgids, nil
}
