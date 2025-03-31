package backend

import (
	"database/sql"
	"os"
	"testing"
)

func TestUpdateOrganizationName(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		t.Fatalf("failed to enable foreign keys: %v", err)
	}

	// Load schema from init.sql
	initSQL, err := os.ReadFile("../sql/init.sql")
	if err != nil {
		t.Fatalf("failed to read init.sql: %v", err)
	}
	if _, err := db.Exec(string(initSQL)); err != nil {
		t.Fatalf("failed to execute init.sql: %v", err)
	}

	// Insert test data with proper permissions
	_, err = db.Exec(`
    INSERT INTO ORGANIZATION (id, name) VALUES (1, 'TechCorp');
    INSERT INTO ORGANIZATION (id, name) VALUES (2, 'ExistingOrg');
    
    INSERT INTO USER (id, email, password, name, verified) 
    VALUES (1, 'user1@example.com', 'pass', 'User One', 1);
    
    INSERT INTO SESSION (id, userid, expires) 
    VALUES (1, 1, '2030-01-01 00:00:00');
    
    -- Add organization role with all required permissions
    INSERT INTO ORG_ROLE (id, orgid, name, can_read, can_write, can_exec) 
    VALUES (1, 1, 'Admin', 1, 1, 1);
    
    -- Add user to organization members
    INSERT INTO ORG_MEMBER (id, userid, orgid) 
    VALUES (1, 1, 1);
    
    -- Assign admin role to member
    INSERT INTO ORG_MEMBER_ROLE (memberid, roleid) 
    VALUES (1, 1);
`)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Test successful name update
	t.Run("successful update", func(t *testing.T) {
		err := UpdateOrganizationName(db, 1, 1, "New TechCorp")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var orgName string
		err = db.QueryRow("SELECT name FROM ORGANIZATION WHERE id = 1").Scan(&orgName)
		if err != nil {
			t.Fatalf("failed to verify update: %v", err)
		}
		if orgName != "New TechCorp" {
			t.Errorf("expected name 'New TechCorp', got '%s'", orgName)
		}
	})

	// Add additional test cases here...
}
