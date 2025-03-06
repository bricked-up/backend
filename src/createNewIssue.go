package backend

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

func CreateNewIssue(issueid int, title string, description string, tagid int, priorityid int, completed time.Time, cost int, date time.Time, ID int) (int64, error) {

	db, err := sql.Open("sqlite", "bricked-up_prod.db")
	if err != nil {
		return -1, err
	}
	defer db.Close()

	issue, err := db.Exec("INSERT INTO issues (title, description, tagid, priority, created, completed, cost, issueid) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", title, description, tagid, priorityid, date, completed, cost, ID)
	if err != nil {
		return -1, err
	}
	id, err := issue.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}
