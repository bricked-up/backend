package backend

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

const dbPath = "Database/backend/sql/BrickedUpDatabase.sql"

// CreateOrganization creates a new organization and assigns the user to it.
// It takes sessionID, orgName, and userID as parameters instead of extracting them from the request.
func CreateOrganization(sessionID, orgName string, userID int) (int, error) {
	// Validate the inputs
	if sessionID == "" || orgName == "" {
		return 0, fmt.Errorf("missing sessionID or orgName")
	}

	// Open the database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	defer db.Close()

	// Enable foreign keys
	_, _ = db.Exec("PRAGMA foreign_keys = ON;")

	// Validate session ID and check if the user ID exists (Assuming a 'sessions' table exists)
	var storedUserID int
	err = db.QueryRow("SELECT user_id FROM sessions WHERE session_id = ?", sessionID).Scan(&storedUserID)
	if err != nil {
		return 0, fmt.Errorf("invalid session: %v", err)
	}

	// Ensure the provided userID matches the stored userID from session
	if userID != storedUserID {
		return 0, fmt.Errorf("user ID does not match the session ID")
	}

	// Insert new organization and get its ID
	var orgID int
	err = db.QueryRow("INSERT INTO organization(name) VALUES(?) RETURNING id", orgName).Scan(&orgID)
	if err != nil {
		return 0, fmt.Errorf("organization name already exists or failed to insert: %v", err)
	}

	// Create admin role for the organization
	_, err = db.Exec("INSERT INTO organization_role (organization_id, name, can_read, can_write, can_execute) VALUES (?, 'admin', 1, 1, 1)", orgID)
	if err != nil {
		return 0, fmt.Errorf("failed to create admin role: %v", err)
	}

	// Assign user to organization
	_, err = db.Exec("INSERT INTO organization_member (user_id, organization_id) VALUES (?, ?)", userID, orgID)
	if err != nil {
		return 0, fmt.Errorf("failed to add user to organization: %v", err)
	}

	// Return the ID of the newly created organization
	return orgID, nil
}
