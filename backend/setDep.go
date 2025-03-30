package backend

import (
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

func SetDep(db *sql.DB, issueAid int, issueBid int, userid int) error {
	// Check if the issues exist
	var existsA, existsB bool

	err := db.QueryRow("SELECT COUNT(*) > 0 FROM ISSUE WHERE id = ?", issueAid).Scan(&existsA)
	if err != nil {
		return err
	}
	if !existsA {
		return errors.New("issueAid does not exist")
	}

	err = db.QueryRow("SELECT COUNT(*) > 0 FROM ISSUE WHERE id = ?", issueBid).Scan(&existsB)
	if err != nil {
		return err
	}
	if !existsB {
		return errors.New("issueBid does not exist")
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
	rows, err := db.Query(query, issueBid, issueAid, userid)
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
	if !access[issueAid] || !access[issueBid] {
		return errors.New("user lacks write access to one or both issues")
	}

	// Check for existing dependency
	var alreadyExists bool
	err = db.QueryRow("SELECT COUNT(*) > 0 FROM DEPENDENCY WHERE issueid = ? AND dependency = ?", issueBid, issueAid).Scan(&alreadyExists)
	if err != nil {
		return err
	}
	if alreadyExists {
		return errors.New("dependency already exists")
	}

	// Insert the dependency
	_, err = db.Exec("INSERT INTO DEPENDENCY (issueid, dependency) VALUES (?, ?)", issueBid, issueAid)
	if err != nil {
		return err
	}

	return nil
}
