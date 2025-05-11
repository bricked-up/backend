package issues

import (
	"brickedup/backend/utils"
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

// CreateIssue creates a new issue in the database with the given parameters.
func CreateIssue(
	sessionid int,
	projectid int,
	title string, 
	desc string, 
	tagid int, 
	priority int, 
	completed time.Time, 
	cost int, 
	date time.Time,
	assignee int, 
	db *sql.DB) (int64, error) {

	var userID int
	var sessionExpires time.Time
	err := db.QueryRow(`
		SELECT userid, expires FROM SESSION 
		WHERE id = ? AND expires > ?
	`, sessionid, time.Now()).Scan(&userID, &sessionExpires)

	if err != nil {
		return -1, err
	}

	// var hasWritePrivilege bool
	// err = db.QueryRow(`
	// 	SELECT EXISTS (
	// 		SELECT 1 FROM PROJECT_MEMBER pm
	// 		JOIN PROJECT_MEMBER_ROLE pmr ON pm.id = pmr.memberid
	// 		JOIN PROJECT_ROLE pr ON pmr.roleid = pr.id
	// 		WHERE pm.userid = ? AND pm.projectid = ? AND pr.can_write = 1
	// 	)
	// `, userID, projectid).Scan(&hasWritePrivilege)
	//
	// if err != nil {
	// 	return -1, err
	// }
	//
	// if !hasWritePrivilege {
	// 	return -1, sql.ErrNoRows // Indicates no matching privileges found
	// }
	//
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

	if projectid >= 0  {
		_, err = db.Exec(`
		INSERT INTO PROJECT_ISSUES (projectid, issueid)
		VALUES (?, ?)
		`, projectid, id)

		if err != nil {
			return -1, err
		}
	}

	if assignee >= 0 {
		var exists bool
	 err = db.QueryRow(`
	 SELECT EXISTS (
		SELECT * FROM project_member WHERE userid = ?
	 )
	 `, assignee).Scan(&exists)
	 if err != nil {
		return -1, err
	 }
	}
		
	_, err  = db.Exec(`
	INSERT INTO USER_ISSUES (userid, issueid)
	VALUES (?, ?)
	`, userID, id)

	if err != nil {
		return -1, err
	}

	return id, nil
}
