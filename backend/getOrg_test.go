package backend

import (
	"database/sql"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

func TestGetOrg(t *testing.T) {

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	defer db.Close()

	initSQL, err := os.ReadFile("../sql/init.sql")
	if err != nil {
		t.Fatalf("failed to open init.sql:%v", err)
	}

	if _, err := db.Exec(string(initSQL)); err != nil {
		t.Fatalf("failed to execute init.sql: %v", err)
	}

	populateSQL, err := os.ReadFile("../sql/populate.sql")

	if err != nil {
		t.Fatalf("failed to open populate.sql:%v", err)
	}

	if _, err := db.Exec(string(populateSQL)); err != nil {
		t.Fatalf("failed tp execute the populate.sql:%v", err)

	}

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
			wantJSON: `{"id":1,"name":"TechCorp Solutions"}`,
		},
		{
			id:       2,
			name:     "Creative Designs Inc",
			wantErr:  false,
			wantJSON: `{"id":2,"name":"Creative Designs Inc"}`,
		},
		{
			id:       3,
			name:     "Data Innovations LLC",
			wantErr:  false,
			wantJSON: `{"id":3,"name":"Data Innovations LLC"}`,
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
