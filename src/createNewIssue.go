package backend

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

// CreateNewIssue creates a new issue in the database with the given parameters.
func CreateNewIssue(title string, desc string, tagid int, priorityid int, completed time.Time, cost int, date time.Time, db *sql.DB) (int64, error) {
	title = sanitizeText(title, TEXT)
	desc = sanitizeText(desc, TEXT)
	issue, err := db.Exec(
		`INSERT INTO issue (title, "desc", tagid, priorityid, created, completed, cost) 
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		title, desc, tagid, priorityid, date, completed, cost,
	)
	if err != nil {
		return -1, err
	}

	id, err := issue.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}
