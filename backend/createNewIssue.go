package backend

import (
	"brickedup/backend/utils"
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

// CreatexNewIssue creates a new issue in the database with the given parameters.
func CreateIssue(title string, desc string, tagid int, priority int, completed time.Time, cost int, date time.Time, db *sql.DB) (int64, error) {
	title = utils.SanitizeText(title, utils.TEXT)
	desc = utils.SanitizeText(desc, utils.TEXT)
	issue, err := db.Exec(
		`INSERT INTO issue (title, "desc", tagid, priority, created, completed, cost) 
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		title, desc, tagid, priority, date, completed, cost,
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
