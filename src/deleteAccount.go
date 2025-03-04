package backend

import (
	"database/sql"

	_ "modernc.org/sqlite" // SQLite driver for database/sql
)

// DeleteUser deletes a user from the dynamically specified table
func deleteUser(sessionid string) error {
	db, err := sql.Open("sqlite", "bricked-up_prod.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM users WHERE id = ? LIMIT 1", sessionid)

	if err != nil {
		return err
	}

	return nil
}
