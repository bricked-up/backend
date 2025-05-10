package projects_test

import (
	"brickedup/backend/projects"
	"brickedup/backend/utils"
	"testing"
	"time"
)

func TestRemoveProjMember(t *testing.T) {
	db := utils.SetupTest(t)
	defer db.Close()

	res, err := db.Exec(`
		INSERT INTO SESSION(userid, expires)
		VALUES(1, ?)
	`, time.Now().Add(24 * time.Hour))

	if err != nil {
		t.Fatal(err.Error())
	}

	session, err := res.LastInsertId()
	if err != nil {
		t.Fatal(err.Error())
	}


	tests := []struct {
		name string // description of this test case
		sessionid int64
		memberid  int
		wantErr   bool
	}{
		{
			name: "Successful",
			sessionid: session,
			memberid: 1,
			wantErr: false,
		},
		{
			name: "Invalid Member",
			sessionid: session,
			memberid: 99,
			wantErr: true,
		},
		{
			name: "Invalid Session",
			sessionid: 0,
			memberid: 99,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := projects.RemoveProjMember(db, tt.sessionid, tt.memberid)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("RemoveProjMember() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("RemoveProjMember() succeeded unexpectedly")
			}
		})
	}
}

