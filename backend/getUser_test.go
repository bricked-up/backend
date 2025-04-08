package backend

import (
	"brickedup/backend/utils"
	"testing"
	"encoding/json"

	_ "modernc.org/sqlite"
)


// TestGetUser demonstrates using an in-memory DB to test GetUser
func TestGetUser(t *testing.T) {
	db := utils.SetupTest(t)

	validUser, err := json.Marshal(
				User{
					Name: "John Doe",
					Email: "john.doe@example.com",
					Password: "",
					Avatar: "avatar1.png",
					Verified: true,
					Projects: []int{ 1, 2 },
					Organizations: nil,
				})

	if err != nil {
		t.Fatal(err)
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
            name:    "Valid User",
            userID:  1,
            wantErr: false,
            wantJSON: string(validUser),
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
            got, err := GetUser(db, tc.userID)
            
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
