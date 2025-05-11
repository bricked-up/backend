package projects

import (
	"brickedup/backend/utils"
	"reflect"
	"testing"
)

func TestGetAllUsers(t *testing.T) {
	db := utils.SetupTest(t)
	defer db.Close()

	expected := []int{1, 2, 3, 4, 5, 6}

	got, err := GetAllProj(db)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, got) {
		t.Fatal("Expected:", expected, ", got:", got)
	}
}
