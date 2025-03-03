package backend

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

// db is a global variable that holds the database connection.
// It is initialized in the `init` function and used throughout the package.
var db *sql.DB

// init initializes the database connection when the program starts.
// This ensures that the database is ready to use when the package is imported.
func init() {
	var err error
	//Open a connection to the SQLite database".
	db, err = sql.Open("sqlite", "./bricked-up_prod.db")
	if err != nil {
		//If the connection fails, log the error and stop the program.
		//This ensures the application does not run with a broken database connection.
		log.Fatal("Database connection failed:", err)
	}
}

// assignOrgRoleToUser assigns a role to a user in an organization.
// It takes the following parameters:
//   - userID: The ID of the user to whom the role will be assigned.
//   - orgID: The ID of the organization where the role will be assigned.
//   - roleID: The ID of the role to assign to the user.
//
// It returns an error if the assignment fails.
func assignOrgRoleToUser(userID, orgID, roleID int) error {
	//Check if the user exists in the USER table.
	var userExists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM USER WHERE id = ?)", userID).Scan(&userExists)
	if err != nil {
		//If the query fails, return an error with details.
		return fmt.Errorf("failed to check user existence: %v", err)
	}
	if !userExists {
		//If the user does not exist, return an error.
		return fmt.Errorf("user with ID %d does not exist", userID)
	}

	//Check if the organization exists in the ORGANIZATION table.
	var orgExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM ORGANIZATION WHERE id = ?)", orgID).Scan(&orgExists)
	if err != nil {
		//If the query fails, return an error with details.
		return fmt.Errorf("failed to check organization existence: %v", err)
	}
	if !orgExists {
		//If the organization does not exist, return an error.
		return fmt.Errorf("organization with ID %d does not exist", orgID)
	}

	//Check if the role exists in the ORG_ROLE table for the given organization.
	var roleExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM ORG_ROLE WHERE id = ? AND orgid = ?)", roleID, orgID).Scan(&roleExists)
	if err != nil {
		// If the query fails, return an error with details.
		return fmt.Errorf("failed to check role existence: %v", err)
	}
	if !roleExists {
		// If the role does not exist in the organization, return an error.
		return fmt.Errorf("role with ID %d does not exist in organization %d", roleID, orgID)
	}

	//Check if the user is already a member of the organization.
	var isMember bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM ORG_MEMBER WHERE userid = ? AND orgid = ?)", userID, orgID).Scan(&isMember)
	if err != nil {
		//If the query fails, return an error with details.
		return fmt.Errorf("failed to check organization membership: %v", err)
	}

	if isMember {
		//If the user is already a member, update their role in the ORG_MEMBER table.
		_, err = db.Exec("UPDATE ORG_MEMBER SET roleid = ? WHERE userid = ? AND orgid = ?", roleID, userID, orgID)
		if err != nil {
			// If the update fails, return an error with details.
			return fmt.Errorf("failed to update user role: %v", err)
		}
	} else {
		//If the user is not a member, add them to the organization with the specified role.
		_, err = db.Exec("INSERT INTO ORG_MEMBER (userid, orgid, roleid) VALUES (?, ?, ?)", userID, orgID, roleID)
		if err != nil {
			// If the insert fails, return an error with details.
			return fmt.Errorf("failed to assign role to user: %v", err)
		}
	}

	//If everything succeeds, return nil (no error).
	return nil
}
