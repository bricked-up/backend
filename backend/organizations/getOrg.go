package organizations

import (
	"brickedup/backend/utils"
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

// GetOrgMembers retrieves all member IDs belonging to a specific organization.
func getOrgMembers(db *sql.DB, org *utils.Organization) error {
	org.Members = nil

	// Prepare query to get all members of the organization
	query := `
		SELECT userid 
		FROM ORG_MEMBER 
		WHERE orgid = ?
	`

	// Execute the query
	rows, err := db.Query(query, org.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Collect all user IDs
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return err
		}
		org.Members = append(org.Members, userID)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

// GetOrgProjects fetches all projects belonging to an organization.
func getOrgProjects(db *sql.DB, org *utils.Organization) error {
	org.Projects = nil

	rows, err := db.Query(
		`SELECT id
		FROM PROJECT 
		WHERE orgid = ?`, 
		org.ID)

	if err != nil {
		return err
	}

	for rows.Next() {
		var projid int

		err = rows.Scan(&projid)
		if err != nil {
			return err
		}

		org.Projects = append(org.Projects, projid)
	}

	return nil
}

// GetOrgRoles retrieves an array of all roles belonging to an organization.
func getOrgRoles(db *sql.DB, org *utils.Organization) error {
	org.Roles = nil

	rows, err := db.Query(
		`SELECT id
		FROM ORG_ROLE
		WHERE orgid = ?`, 
		org.ID)

	if err != nil {
		return err
	}

	for rows.Next() {
		var roleid int

		err = rows.Scan(&roleid)
		if err != nil {
			return err
		}

		org.Roles = append(org.Roles, roleid)
	}

	return nil
}


// GetOrg returns an organization entry.
func GetOrg(db *sql.DB, orgid int) (*utils.Organization, error) {
	row := db.QueryRow(`SELECT id, name FROM organization where id = ?`, orgid)

	org := &utils.Organization{}
	if err := row.Scan(&org.ID, &org.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Organization not found")
		}
		return nil, err
	}

	if err := getOrgMembers(db, org); err != nil {
		return nil, err
	}

	if err := getOrgProjects(db, org); err != nil {
		return nil, err
	}

	if err := getOrgRoles(db, org); err != nil {
		return nil, err
	}

	return org, nil
}
