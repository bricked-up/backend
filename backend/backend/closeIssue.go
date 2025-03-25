package backend

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	_ "modernc.org/sqlite"
)

func CloseIssue(db *sql.DB, issueid string, sessionid string) error {
	// Convert IDs to integers
	intSessionID, err := strconv.Atoi(sessionid)
	if err != nil {
		return fmt.Errorf("invalid session ID: %v", sessionid)
	}
	intIssueID, err := strconv.Atoi(issueid)
	if err != nil {
		return fmt.Errorf("invalid issue ID: %v", issueid)
	}

	// Step 1: Get user ID from session
	var userID int
	err = db.QueryRow("SELECT userid FROM session WHERE id = ?", intSessionID).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no session found for session ID %v", intSessionID)
		}
		return err
	}

	// Step 2: Check permission to close issue
	var canWrite bool
	query := `
		SELECT pr.can_write
		FROM ISSUE i
		JOIN PROJECT_ISSUES pi ON i.id = pi.issueid
		JOIN PROJECT_MEMBER pm ON pi.projectid = pm.projectid
		JOIN PROJECT_MEMBER_ROLE pmr ON pm.id = pmr.memberid
		JOIN PROJECT_ROLE pr ON pmr.roleid = pr.id
		WHERE i.id = ? AND pm.userid = ?;
	`
	err = db.QueryRow(query, intIssueID, userID).Scan(&canWrite)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("no matching session or permission, try again")
		}
		return err
	}

	if !canWrite {
		return errors.New("user does not have privileges to close the issue")
	}

	// Step 3: Delete the issue
	_, err = db.Exec(`DELETE FROM issue WHERE id = ?`, intIssueID)
	if err != nil {
		return err
	}

	fmt.Println("Issue closed successfully")
	return nil
}
