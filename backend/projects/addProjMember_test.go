package projects

import (
	"brickedup/backend/utils"
	"testing"
	"time"
)

func TestAddProjMember(t *testing.T) {
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
		userid    int
		roleid    int
		projectid int
		wantErr   bool
	}{
		{
			name: "Successful",
			sessionid: session,
			userid: 1,
			roleid: 1,
			projectid: 1,
			wantErr: false,
		},
		{
			name: "Inexistant Project",
			sessionid: session,
			userid: 1,
			roleid: 1,
			projectid: 1000,
			wantErr: true,
		},
		{
			name: "Inexistant Role",
			sessionid: session,
			userid: 1,
			roleid: 1000,
			projectid: 1,
			wantErr: true,
		},
		{
			name: "Inexistant User",
			sessionid: session,
			userid: 1000,
			roleid: 1,
			projectid: 1,
			wantErr: true,
		},
		{
			name: "Inexistant Project Manager",
			sessionid: 1000,
			userid: 1,
			roleid: 1,
			projectid: 1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := AddProjMember(db, tt.sessionid, tt.userid, tt.roleid, tt.projectid)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("AddProjMember() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("AddProjMember() succeeded unexpectedly")
			}
		})
	}
}

