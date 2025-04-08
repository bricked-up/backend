package issues

import (
	"brickedup/backend/utils"
	"testing"

	_ "modernc.org/sqlite"
)

func TestGetDep(t *testing.T) {
	db := utils.SetupTest(t)

	// Define subtests
	tests := []struct {
		name     string
		issueID  int
		wantErr  bool
		wantJSON string
	}{
		{
			name:     "Issue #3 has dependency 1",
			issueID:  3,
			wantErr:  false,
			wantJSON: `[1]`,
		},
		{
			name:     "Issue #4 has dependency 2",
			issueID:  4,
			wantErr:  false,
			wantJSON: `[2]`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := getDep(db, tc.issueID)

			if tc.wantErr && err == nil {
				t.Fatalf("expected an error but got none")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("did not expect an error but got: %v", err)
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
