package backend

import (
	"database/sql"
	"encoding/json"
	"errors"

	_ "modernc.org/sqlite"
)

func getUserProjects(db *sql.DB, userid int, user *User) error {
	rows, err := db.Query(
		`SELECT projectid 
		FROM PROJECT_MEMBER
		WHERE userid = ?`, userid)

	if err != nil {
		return err
	}

	for rows.Next() {
		var projectid int

		err := rows.Scan(&projectid)
		if err != nil {
			return err
		}

		user.Projects = append(user.Projects, projectid)
	}

	return nil
}

func getUserOrganizations(db *sql.DB, userid int, user *User) error {
	rows, err := db.Query(
		`SELECT orgid 
		FROM ORG_MEMBER
		WHERE userid = ?`, userid)

	if err != nil {
		return err
	}

	for rows.Next() {
		var orgid int

		err := rows.Scan(&orgid)
		if err != nil {
			return err
		}

		user.Organizations = append(user.Organizations, orgid)
	}

	return nil
}


// GetUser fetches one user by sessionid from the DB and returns JSON data.
func GetUser(db *sql.DB, sessionid int) ([]byte, error) {

	// Get the user ID from the session
	var userID int
	err := db.QueryRow("SELECT userid FROM SESSION WHERE id = ?", sessionid).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no session exists for the provided sessionID")
		}
		return nil, err
	}

	// Get exactly one row for the given userID.
	row := db.QueryRow(`SELECT name, email, verified, avatar FROM USER WHERE id = ?`, userID)

	var user User
	// Scan fills our user struct with the row's data or returns an error if no row.
	err = row.Scan(&user.Name, &user.Email, &user.Verified, &user.Avatar)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("UserId not found")
		}
		return nil, err
	}

	err = getUserProjects(db, userID, &user)
	if err != nil {
		return nil, err
	}

	// Convert the user struct to JSON.
	jsonUser, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	return jsonUser, nil
}
