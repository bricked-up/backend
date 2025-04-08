package users

import (
	"brickedup/backend/utils"
	"testing"

	_ "modernc.org/sqlite"
)

// TestUpdateUser demonstrates using an in-memory DB to test UpdateUser.
func TestUpdateUser(t *testing.T) {
	db := utils.SetupTest(t)
	defer db.Close()

	var originalUser utils.User
	err := db.QueryRow(
		` SELECT name, email, password, avatar
		FROM USER 
		WHERE id = 1`).Scan(
			&originalUser.Name,
			&originalUser.Email,
			&originalUser.Password,
			&originalUser.Avatar,
		)

    if err != nil {
        t.Errorf("failed to get user: %v", err)
    }

	updatedUser := originalUser
	updatedUser.Name = "Ivan123"

    err = UpdateUser(db, 1, &updatedUser)
    if err != nil {
        t.Errorf("ChangeDisplayName returned error: %v", err)
    }

    var updatedName string
    err = db.QueryRow("SELECT name FROM USER WHERE id = 1").Scan(&updatedName)
    if err != nil {
        t.Errorf("failed to query updated name: %v", err)
    }
    if updatedName == originalUser.Name {
        t.Errorf("name was not changed from '%s' to '%s'", originalUser.Name, updatedName)
    }
}
