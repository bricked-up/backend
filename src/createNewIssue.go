package backend

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

// CreateNewIssue creates a new issue in the database
// inputs arethe title, description, tagid, priorityid, completed, cost, date, and ID
// returns are the id and error
func CreateNewIssue(title string, description string, tagid int, priorityid int, completed time.Time, cost int, date time.Time, ID int) (int64, error) {
	//Open the database
	db, err := sql.Open("sqlite", "bricked-up_prod.db")
	if err != nil {
		return -1, err
	}
	defer db.Close()

	//Insert the new issue into the database
	issue, err := db.Exec("INSERT INTO issues (title, tagid, priority, created, completed, cost, issueid) VALUES (?, ?, ?, ?, ?, ?, ?)", title, tagid, priorityid, date, completed, cost, ID)
	if err != nil {
		return -1, err
	}
	// Get the id of the new issue
	id, err := issue.LastInsertId()
	if err != nil {
		return -1, err
	}
	//Return the id and nil
	return id, nil
}
