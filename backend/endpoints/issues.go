package endpoints

import (
	"brickedup/backend/issues"
	"brickedup/backend/utils"
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"
)

// CreateIssueHandler handles POST requests to insert a new issue into the database.
func CreateIssueHandler (db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w,"Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	session := r.FormValue("sessionid")
	sessionID, err := strconv.Atoi(session)
	if err != nil {
		http.Error(w,"Invalid session ID", http.StatusBadRequest)
		return
	}

	// Extract and validate form values
	projectid, err := strconv.Atoi(r.FormValue("projectid"))
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	desc := r.FormValue("desc")
	tagID, err := strconv.Atoi(r.FormValue("tagid"))
	if err != nil {
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	priority, err := strconv.Atoi(r.FormValue("priority"))
	if err != nil {
		http.Error(w, "Invalid priority ID", http.StatusBadRequest)
		return
	}

	cost, err := strconv.Atoi(r.FormValue("cost"))
	if err != nil {
		http.Error(w, "Invalid cost", http.StatusBadRequest)
		return
	}

	assignee, err := strconv.Atoi(r.FormValue("assignee"))
	if err != nil {
		http.Error(w, "Invalid assignee", http.StatusBadRequest)
		return
	}

	// Call your backend logic
	_, err = issues.CreateIssue(
		sessionID, 
		projectid, 
		title, 
		desc, 
		tagID, 
		priority, 
		cost, 
		time.Now(),
		assignee, 
		db)

	if err != nil {
		http.Error(w, "Failed to create issue: "+err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
}

// GetIssueHandler handles GET requests and retrieves information
// on the issue on /get-issue.
// The `issueid` is specified as a URL parameter.
func GetIssueHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameter
	issueIDStr := r.URL.Query().Get("issueid")
	if issueIDStr == "" {
		http.Error(w, "Missing issueid parameter", http.StatusBadRequest)
		return
	}

	issueID, err := strconv.Atoi(issueIDStr)
	if err != nil {
		http.Error(w, "Invalid issueid", http.StatusBadRequest)
		return
	}

	// Fetch issue details
	jsonStr, err := issues.GetIssue(db, issueID)
	if err == sql.ErrNoRows {
		http.Error(w, "Issue not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to fetch issue: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonStr))
}

// UpdateIssueHandler handles PATCH requests to update an existing issue on /update-issue.
func UpdateIssueHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPatch {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    if err := r.ParseForm(); err != nil {
        http.Error(w, "Invalid form", http.StatusBadRequest)
        return
    }

    // Parse and validate issue ID
    issueIDStr := r.FormValue("issueid")
    issueID, err := strconv.Atoi(issueIDStr)
    if err != nil {
        http.Error(w, "Invalid issue ID", http.StatusBadRequest)
        return
    }

    // Build the payload
    var issue utils.Issue
    issue.Title = r.FormValue("title")
    issue.Desc  = r.FormValue("desc")

    if costStr := r.FormValue("cost"); costStr != "" {
        cost, err := strconv.Atoi(costStr)
        if err != nil {
            http.Error(w, "Invalid cost", http.StatusBadRequest)
            return
        }
        issue.Cost = cost
    }

    if tagStr := r.FormValue("tagid"); tagStr != "" {
        tagID, err := strconv.Atoi(tagStr)
        if err != nil {
            http.Error(w, "Invalid tag ID", http.StatusBadRequest)
            return
        }
        issue.TagID = tagID
    }

    if prioStr := r.FormValue("priority"); prioStr != "" {
        priority, err := strconv.Atoi(prioStr)
        if err != nil {
            http.Error(w, "Invalid priority", http.StatusBadRequest)
            return
        }
        issue.Priority = priority
    }

    // Parse optional completed timestamp
    completedStr := r.FormValue("completed")
    if completedStr != "" {
        const layout = "2006-01-02 15:04:05"
        t, err := time.Parse(layout, completedStr)
        if err != nil {
            http.Error(w, "Invalid completed date format (use YYYY-MM-DD HH:MM:SS)", http.StatusBadRequest)
            return
        }
        issue.Completed = sql.NullTime{Time: t, Valid: true}
    }

    // Call core business logic
    if err := issues.UpdateIssue(db, issueID, &issue); err != nil {
        if err.Error() == "no issue found for issue ID "+issueIDStr {
            http.Error(w, "Issue not found", http.StatusNotFound)
        } else {
            log.Println("UpdateIssue error:", err)
            http.Error(w, "Failed to update issue: "+err.Error(), http.StatusInternalServerError)
        }
        return
    }

    w.WriteHeader(http.StatusOK)
}
