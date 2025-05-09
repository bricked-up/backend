package projects

import (
	"brickedup/backend/utils"
	"reflect"
	"testing"

	_ "modernc.org/sqlite"
)

func TestGetProjMember(t *testing.T) {
	db := utils.SetupTest(t)
	defer db.Close()

	validUser := utils.ProjectMember{
					ID: 1,
					UserID: 1,
					ProjectID: 1,
					Roles: []int{1},
					CanExec: true,
					CanWrite: true,
					CanRead: true,
					Issues: []int{1},
				}

    // Definition of subtests that check if part of the data from the json is 
	// corresponding to the actual db.
    tests := []struct {
        name     string
        userID   int
        wantErr  bool
        want	 *utils.ProjectMember
    }{
        {
            name:    "Valid User",
            userID:  1,
            wantErr: false,
            want: &validUser,
        },
        {
            // Non-existent row = expect an error
            name:    "User #999 - Does not exist",
            userID:  999,
            wantErr: true,
            want: nil,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            got, err := GetProjMember(db, tc.userID)
            
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

            if !reflect.DeepEqual(got, tc.want) {
                t.Errorf("Error mismatch.\nGot: %+v\nWant: %+v", got, tc.want)
            }
        })
    }
}
