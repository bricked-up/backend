package backend

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

const dbPath = "Database/backend/sql/BrickedUpDatabase.sql"

func createOrganization(w http.ResponseWriter, r *http.Request) {
	sessionID := r.FormValue("sessionid")
	orgName := r.FormValue("name")

	if sessionID == "" || orgName == "" {
		http.Error(w, "Missing sessionid or name", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Enable foreign keys
	_, _ = db.Exec("PRAGMA foreign_keys = ON;")

	// Validate session ID and get user ID (Assuming a 'sessions' table exists)
	var userID int
	err = db.QueryRow("SELECT user_id FROM sessions WHERE session_id = ?", sessionID).Scan(&userID)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	// Insert new organization and get its ID
	var orgID int
	err = db.QueryRow("INSERT INTO organization(name) VALUES(?) RETURNING id", orgName).Scan(&orgID)
	if err != nil {
		http.Error(w, "Organization name already exists", http.StatusConflict)
		return
	}

	// Create admin role for the organization
	_, err = db.Exec("INSERT INTO organization_role (organization_id, name, can_read, can_write, can_execute) VALUES (?, 'admin', 1, 1, 1)", orgID)
	if err != nil {
		http.Error(w, "Failed to create admin role", http.StatusInternalServerError)
		return
	}

	// Assign user to organization
	_, err = db.Exec("INSERT INTO organization_member (user_id, organization_id) VALUES (?, ?)", userID, orgID)
	if err != nil {
		http.Error(w, "Failed to add user to organization", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Organization created successfully with ID %d", orgID)
}

func main() {
	http.HandleFunc("/createOrganization", createOrganization)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
