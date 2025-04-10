package endpoints

import (
	"brickedup/backend/projects"
	"database/sql"
	"log"
	"net/http"
	"strconv"
)

// CreateTagHandler handles POST requests to create a new tag associated with a
// project on /create-tag.
// It validates the session, form inputs, and calls the CreateTag logic function.
func CreateTagHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Read session ID from cookie
	cookie, err := r.Cookie(LoginCookie)
	if err != nil {
		http.Error(w, "Missing session cookie", http.StatusUnauthorized)
		return
	}
	sessionID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	// Parse form values
	projectID, err := strconv.Atoi(r.FormValue("projectid"))
	if err != nil {
		http.Error(w, "Invalid or missing project ID", http.StatusBadRequest)
		return
	}

	tagName := r.FormValue("name")
	tagColor := r.FormValue("color")

	// Call core logic
	_, err = projects.CreateTag(db, sessionID, projectID, tagName, tagColor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("createTag error:", err)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
}

// DeleteTagHandler handles DELETE requests to delete a tag by its ID on /delete-tag.
// The session is used to validate whether the user has permission to delete the tag.
func DeleteTagHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form values
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Retrieve session ID from cookie
	cookie, err := r.Cookie(LoginCookie)
	if err != nil {
		http.Error(w, "Missing session cookie", http.StatusUnauthorized)
		return
	}
	sessionID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	// Parse tag ID from form
	tagIDStr := r.FormValue("tagid")
	if tagIDStr == "" {
		http.Error(w, "Missing tag ID", http.StatusBadRequest)
		return
	}

	tagID, err := strconv.Atoi(tagIDStr)
	if err != nil {
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	// Call core logic
	err = projects.DeleteTag(db, sessionID, tagID)
	if err != nil {
		http.Error(w, "Failed to delete tag: "+err.Error(), http.StatusForbidden)
		log.Println("deleteTag error:", err)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
}
