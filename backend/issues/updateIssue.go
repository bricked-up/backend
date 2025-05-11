package issues

import (
    "brickedup/backend/utils"
    "database/sql"
    "errors"
    "strconv"

    _ "modernc.org/sqlite"
)

// UpdateIssue retrieves the issue by ID, sanitizes its Title & Desc,
// and writes the new values back to the ISSUE table in one shot (autocommit).
func UpdateIssue(db *sql.DB, issueID int, issue *utils.Issue) error {
    // 1) Sanitize free-text fields
    issue.Title = utils.SanitizeText(issue.Title, utils.TEXT)
    issue.Desc  = utils.SanitizeText(issue.Desc,  utils.TEXT)

    // 2) Verify the issue exists
    var id int
    err := db.QueryRow(
        "SELECT id FROM ISSUE WHERE id = ?",
        issueID,
    ).Scan(&id)
    if err != nil {
        if err == sql.ErrNoRows {
            return errors.New("no issue found for issue ID " + strconv.Itoa(issueID))
        }
        return err
    }

    // 3) Perform the UPDATE (autocommitted)
    const q = `
        UPDATE ISSUE
           SET title     = ?,
               desc      = ?,
               cost      = ?,
               tagid     = ?,
               priority  = ?,
               completed = ?
         WHERE id = ?
    `
    _, err = db.Exec(
        q,
        issue.Title,
        issue.Desc,
        issue.Cost,
        issue.TagID,
        issue.Priority,
        issue.Completed,
        issueID,
    )
    return err
}
