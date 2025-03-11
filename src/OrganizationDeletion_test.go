package backend

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

// TestDeleteOrganization tests the DeleteOrganization function.
func TestDeleteOrganization(t *testing.T) {
	// Open in-memory database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	defer db.Close()

	// Create necessary tables and insert test data
	_, err = db.Exec(`
		CREATE TABLE sessions (session_id TEXT, user_id INTEGER);
		CREATE TABLE organizations (id INTEGER, name TEXT);
		CREATE TABLE organization_roles (organization_id INTEGER, user_id INTEGER, role_id INTEGER);
		CREATE TABLE roles (id INTEGER, can_exec BOOLEAN);

		INSERT INTO sessions (session_id, user_id) VALUES ('valid_session', 1);
		INSERT INTO organizations (id, name) VALUES (1, 'Test Organization');
		INSERT INTO organization_roles (organization_id, user_id, role_id) VALUES (1, 1, 1);
		INSERT INTO roles (id, can_exec) VALUES (1, TRUE);
	`)
	if err != nil {
		t.Fatalf("failed to set up test data: %v", err)
	}

	tests := []struct {
		name      string
		sessionID string
		orgID     string
		wantErr   bool
	}{
		{"ValidSessionAndOrgID", "valid_session", "1", false},
		{"MissingSessionID", "", "1", true},
		{"MissingOrgID", "valid_session", "", true},
		{"InvalidSessionID", "invalid_session", "1", true},
		{"InvalidOrgID", "valid_session", "999", true},
		{"NoPermission", "valid_session", "1", true}, // Assuming user does not have permission
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Pass db connection to the DeleteOrganization function
			err := DeleteOrganization(db, tt.sessionID, tt.orgID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteOrganization() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
