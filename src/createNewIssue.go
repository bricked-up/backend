package backend

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

// CreateNewIssue creates a new issue in the database with the given parameters.
// Input: issueid, title, description, tagid, priorityid, completed, cost, date, ID
// Output: int64, errors
func CreateNewIssue(issueid int, title string, tagid int, priorityid int, completed time.Time, cost int, date time.Time, ID int, db *sql.DB) (int64, error) {

	issue, err := db.Exec("INSERT INTO issues (title, tagid, priority, created, completed, cost, issueid) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", title, tagid, priorityid, date, completed, cost, ID)
	if err != nil {
		return -1, err
	}

	defer db.Close()
	id, err := issue.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}
