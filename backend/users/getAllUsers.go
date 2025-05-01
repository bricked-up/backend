package users

import "database/sql"

// GetAllUsers retrieves the IDs of all the verified users.
func GetAllUsers(db *sql.DB) ([]int, error) {
	var userids []int

	res, err := db.Query(`
		SELECT id FROM USER
		WHERE verified = 1
		`)

	if err != nil {
		return nil, err
	}

	for res.Next() {
		var userid int
		err = res.Scan(&userid)

		if err != nil {
			return nil, err
		}

		userids = append(userids, userid)
	}

	return userids, nil
}
