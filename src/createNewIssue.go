package backend

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

// CreateNewIssue creates a new issue in the database with the given parameters.
// Input: issueid, title, description, tagid, priorityid, completed, cost, date, ID
// Output: int64, errors
func CreateNewIssue(issueid int, desc string, title string, tagid int, priorityid int, completed time.Time, cost int, date time.Time, db *sql.DB) (int64, error) {
	issue, err := db.Exec("INSERT INTO issue (id, title, desc, tagid, priorityid, created, completed, cost) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", issueid, title, desc, tagid, priorityid, date, completed, cost)
	if err != nil {
		return -1, err
	}

	id, err := issue.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}
