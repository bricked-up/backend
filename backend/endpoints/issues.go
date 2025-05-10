package endpoints

import (
	"brickedup/backend/issues"
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

	// Parse date and completed (optional: validate layout)
	dateStr := r.FormValue("date")
	completedStr := r.FormValue("completed")

	const layout = "2006-01-02 15:04:05"

	date, err := time.Parse(layout, dateStr)
	if err != nil {
		http.Error(w, "Invalid date format (use YYYY-MM-DD HH:MM:SS)", http.StatusBadRequest)
		return
	}

	completed, err := time.Parse(layout, completedStr)
	if err != nil {
		http.Error(w, "Invalid completed date format (use YYYY-MM-DD HH:MM:SS)", http.StatusBadRequest)
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
		completed, 
		cost, 
		date, 
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

