package backend

import (
    "database/sql"
    "os"
    "testing"

    _ "modernc.org/sqlite"
)

// TestGetUserDetails demonstrates using an in-memory DB to test getUserDetails
func TestGetUserDetails(t *testing.T) {
    // Open an in-memory DB
    db, err := sql.Open("sqlite", ":memory:")
    if err != nil {
        t.Fatalf("failed to open in-memory db: %v", err)
    }
    defer db.Close()

    // Exec init.sql (create tables, etc.)
    initSQL, err := os.ReadFile("../sql/init.sql")
    if err != nil {
        t.Fatalf("failed to open init.sql: %v", err)
    }
    if _, err := db.Exec(string(initSQL)); err != nil {
        t.Fatalf("failed to exec init.sql: %v", err)
    }

    // Exec populate.sql (insert rows)
    populateSQL, err := os.ReadFile("../sql/populate.sql")
    if err != nil {
        t.Fatalf("failed to open populate.sql: %v", err)
    }
    if _, err := db.Exec(string(populateSQL)); err != nil {
        t.Fatalf("failed to exec populate.sql: %v", err)
    }

    // Definition of subtests that check if part of the data from the json is 
	// corresponding to the actual db.
    tests := []struct {
        name     string
        userID   int
        wantErr  bool
        wantJSON string
    }{
        {
            name:    "User #1 - John Doe",
            userID:  1,
            wantErr: false,
            wantJSON: `{"id":1,"name":"John Doe","email":"john.doe@example.com","password":"","avatar":"avatar1.png","verified":true}`,
        },
        {
            name:    "User #2 - Jane Smith",
            userID:  2,
            wantErr: false,
            wantJSON: `{"id":2,"name":"Jane Smith","email":"jane.smith@example.com","password":"","avatar":"avatar2.png","verified":true}`,
        },
        {
            name:    "User #3 - Mike Johnson",
            userID:  3,
            wantErr: false,
            wantJSON: `{"id":3,"name":"Mike Johnson","email":"mike.johnson@example.com","password":"","avatar":"avatar3.png","verified":true}`,
        },
        {
            name:    "User #4 - Sarah Williams",
            userID:  4,
            wantErr: false,
            wantJSON: `{"id":4,"name":"Sarah Williams","email":"sarah.williams@example.com","password":"","avatar":"avatar4.png","verified":false}`,
        },
        {
            name:    "User #5 - Alex Brown",
            userID:  5,
            wantErr: false,
            wantJSON: `{"id":5,"name":"Alex Brown","email":"alex.brown@example.com","password":"","avatar":"avatar5.png","verified":false}`,
        },
        {
            // Non-existent row = expect an error
            name:    "User #999 - Does not exist",
            userID:  999,
            wantErr: true,
            wantJSON: "",
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            got, err := getUserDetails(db, tc.userID)
            
            if tc.wantErr && err == nil {
                t.Fatalf("expected an error but got none")
            }
            if !tc.wantErr && err != nil {
                t.Fatalf("did not expect an error but got: %v", err)
            }

            // If we expected an error and got one, we can stop here
            if err != nil {
                return
            }

            // Compare the entire JSON string
            gotStr := string(got)
            if gotStr != tc.wantJSON {
                t.Errorf("JSON mismatch.\nGot:  %s\nWant: %s", gotStr, tc.wantJSON)
            }
        })
    }
}
