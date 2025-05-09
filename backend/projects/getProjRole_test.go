package projects

import (
	"brickedup/backend/utils"
	"reflect"
	"testing"
)

func TestGetProjRole(t *testing.T) {
	db := utils.SetupTest(t)
	defer db.Close()

	want := &utils.ProjectRole{
		ID: 1,
		ProjectID: 1,
		Name: "Project Manager",
		CanExec: true,
		CanWrite: true,
		CanRead: true,
	}

	got, err := GetProjRole(db, want.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Error mismatch! Got: %+v, want: %+v!", got, want)
	}

	// Role does not exist.
	got, err = GetProjRole(db, 99999)
	if err == nil {
		t.Fatal("Role that does not exist should return an error!")
	}
}
