package issues

import (
    "database/sql"
    "testing"
    "time"

    "brickedup/backend/utils"
    _ "modernc.org/sqlite"
)

func TestUpdateIssue(t *testing.T) {
    db := utils.SetupTest(t)
    defer db.Close()

    tests := []struct {
        name     string
        seed     bool
        issueID  int
        update   *utils.Issue
        wantErr  bool
        errMsg   string
    }{
        {
            name:    "success",
            seed:    true,
            issueID: 1,
            update: &utils.Issue{
                Title:    "New Title",
                Desc:     "New Desc",
                Cost:     20,
                TagID:    2,
                Priority: 5,
                Completed: sql.NullTime{Valid: false},
            },
            wantErr: false,
        },
        {
            name:    "not found",
            seed:    false,
            issueID: 999,
            update: &utils.Issue{
                Title:    "X",
                Desc:     "Y",
                Cost:     1,
                TagID:    1,
                Priority: 1,
                Completed: sql.NullTime{Valid: false},
            },
            wantErr: true,
            errMsg:  "no issue found for issue ID 999",
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            // Clear any existing rows so our explicit ID insert won't collide
            if _, err := db.Exec(`DELETE FROM ISSUE`); err != nil {
                t.Fatalf("clearing ISSUE table: %v", err)
            }

            if tc.seed {
                _, err := db.Exec(`
                    INSERT INTO ISSUE (id, title, desc, cost, tagid, priority, created, completed)
                    VALUES (?, ?, ?, ?, ?, ?, ?, ?)
                `,
                    tc.issueID,
                    "Old Title",
                    "Old Desc",
                    10,
                    tc.update.TagID,
                    tc.update.Priority,
                    time.Now(),
                    nil,
                )
                if err != nil {
                    t.Fatalf("seeding ISSUE failed: %v", err)
                }
            }

            err := UpdateIssue(db, tc.issueID, tc.update)
            if tc.wantErr {
                if err == nil {
                    t.Fatalf("expected error but got nil")
                }
                if err.Error() != tc.errMsg {
                    t.Errorf("error = %q, want %q", err.Error(), tc.errMsg)
                }
                return
            }
            if err != nil {
                t.Fatalf("UpdateIssue returned error: %v", err)
            }

            var gotTitle, gotDesc string
            var gotCost int
            err = db.QueryRow(
                "SELECT title, desc, cost FROM ISSUE WHERE id = ?",
                tc.issueID,
            ).Scan(&gotTitle, &gotDesc, &gotCost)
            if err != nil {
                t.Fatalf("querying updated issue: %v", err)
            }

            if gotTitle != tc.update.Title {
                t.Errorf("title = %q, want %q", gotTitle, tc.update.Title)
            }
            if gotDesc != tc.update.Desc {
                t.Errorf("desc = %q, want %q", gotDesc, tc.update.Desc)
            }
            if gotCost != tc.update.Cost {
                t.Errorf("cost = %d, want %d", gotCost, tc.update.Cost)
            }
        })
    }
}