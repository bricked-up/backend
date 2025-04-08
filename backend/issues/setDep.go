package issues

import (
	"database/sql"
	"errors"
	"strconv"

	_ "modernc.org/sqlite"
)

// SetDep assigns the dependency to the issue.
func SetDep(db *sql.DB, issueid int, dependency int, sessionid int) error {

	// Look up the userID in the SESSION table.
	var userid int
	err := db.QueryRow("SELECT userid FROM SESSION WHERE id = ?", sessionid).Scan(&userid)
	if err != nil {
		// If no row is found, return a custom error message.
		if err == sql.ErrNoRows {
			return errors.New("no session found for session ID " + strconv.Itoa(sessionid))
		}
		// Otherwise, return the original error from the DB.
		return err
	}

	// Check if the issues exist
	var existsA, existsB bool

	err = db.QueryRow("SELECT COUNT(*) > 0 FROM ISSUE WHERE id = ?", issueid).Scan(&existsA)
	if err != nil {
		return err
	}
	if !existsA {
		return errors.New("issue does not exist: " + strconv.Itoa(issueid))
	}

	err = db.QueryRow("SELECT COUNT(*) > 0 FROM ISSUE WHERE id = ?", dependency).Scan(&existsB)
	if err != nil {
		return err
	}
	if !existsB {
		return errors.New("issue does not exist: " + strconv.Itoa(dependency))
	}

	// Check if the user has write access to both issues
	query := `
        SELECT pi.issueid, pr.can_write
        FROM PROJECT_ISSUES pi
        JOIN PROJECT_MEMBER pm ON pi.projectid = pm.projectid
        JOIN PROJECT_MEMBER_ROLE pmr ON pm.id = pmr.memberid
        JOIN PROJECT_ROLE pr ON pmr.roleid = pr.id
        WHERE pi.issueid IN (?, ?)
          AND pm.userid = ?;
    `
	rows, err := db.Query(query, dependency, issueid, userid)
	if err != nil {
		return err
	}
	defer rows.Close()

	access := map[int]bool{}
	for rows.Next() {
		var issueID int
		var canWrite bool
		if err := rows.Scan(&issueID, &canWrite); err != nil {
			return err
		}
		if canWrite {
			access[issueID] = true
		}
	}

	// Ensure access exists for both issues
	if !access[issueid] || !access[dependency] {
		return errors.New("user lacks write access to one or both issues")
	}

	// Check for existing dependency
	var alreadyExists bool
	err = db.QueryRow(
		`SELECT COUNT(*) > 0 
		FROM DEPENDENCY 
		WHERE issueid = ? AND dependency = ?`, 
		dependency, issueid).Scan(&alreadyExists)

	if err != nil {
		return err
	}
	if alreadyExists {
		return errors.New("dependency already exists")
	}

	// Insert the dependency
	_, err = db.Exec(
		`INSERT INTO DEPENDENCY (issueid, dependency) 
		VALUES (?, ?)`, dependency, issueid)

	if err != nil {
		return err
	}

	return nil
}
