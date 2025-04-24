package organizations

import (
	"brickedup/backend/utils"
	"testing"

	_ "modernc.org/sqlite"
)

func TestGetOrg(t *testing.T) {
	db := utils.SetupTest(t)
	defer db.Close()

	tests := []struct {
		id       int
		name     string
		wantErr  bool
		wantJSON string
	}{
		{
			id:       1,
			name:     "TechCorp Solutions",
			wantErr:  false,
			wantJSON: `{"id":1,"name":"TechCorp Solutions","members":[1,2,3],"projects":[1,2,3],"roles":[1,2,3]}`,
		},
		{
			id:       2,
			name:     "Creative Designs Inc",
			wantErr:  false,
			wantJSON: `{"id":2,"name":"Creative Designs Inc","members":[2,4],"projects":[4,5],"roles":[4,5]}`,
		},
		{
			id:       3,
			name:     "Data Innovations LLC",
			wantErr:  false,
			wantJSON: `{"id":3,"name":"Data Innovations LLC","members":[5,3],"projects":[6],"roles":[6,7]}`,
		},
		{
			id:       999,
			name:     "False Company Inc",
			wantErr:  true,
			wantJSON: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GetOrg(db, tc.id)
			if tc.wantErr && err == nil {
				t.Fatalf("expected an error but got none")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("did not expect an erorr but got: %v", err)
			}

			if err != nil {
				return
			}

			gotStr := string(got)
			if gotStr != tc.wantJSON {
				t.Errorf("JSON mismatch.\nGot:  %s\nWant: %s", gotStr, tc.wantJSON)
			}
		})
	}

}
