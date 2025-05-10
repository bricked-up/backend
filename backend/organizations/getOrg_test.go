package organizations

import (
	"brickedup/backend/utils"
	"reflect"
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
		want 	 *utils.Organization
	}{
		{
			id:       1,
			name:     "TechCorp Solutions",
			wantErr:  false,
			want: 	  &utils.Organization{
				ID: 1,
				Name: "TechCorp Solutions",
				Members: []int{1,2,3},
				Projects: []int{1,2,3},
				Roles: []int{1,2,3},
			},
		},
		{
			id:       999,
			name:     "False Company Inc",
			wantErr:  true,
			want: 	  nil,
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

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Error mismatch.\nGot:  %+v\nWant: %+v", got, tc.want)
			}
		})
	}

}
