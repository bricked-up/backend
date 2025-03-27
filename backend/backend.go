// Package backend provides the backend infrastructure (route handling + database)
// for the Bricked-Up website.
package backend

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type DBHandlerFunc func(*sql.DB, http.ResponseWriter, *http.Request)

// Endpoints maps URL paths to their corresponding handler functions.
var endpoints = map[string]DBHandlerFunc{
	"/login":                   LoginHandler,
    "/signup":                  SignupHandler,
    "/verify":                  VerifyHandler,
    "/change-name":            ChangeNameHandler,
    "/create-issue":           CreateIssueHandler,
    "/create-tag":             CreateTagHandler,
    "/delete-user":            DeleteUserHandler,
    "/delete-tag":             DeleteTagHandler,
    "/issue":                  GetIssueDetailsHandler,
    "/org-members":            GetOrgMembersHandler,
    "/create-org":             CreateOrganizationHandler,
    "/delete-org":             DeleteOrganizationHandler,
    "/remove-org-member-role": RemoveOrgMemberRoleHandler,
    "/assign-org-role":        AssignOrgRoleHandler,
}

// LoginHandler handles requests to the /login endpoint.
// It only allows GET requests and responds with a placeholder message.
func LoginHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method unsupported", http.StatusMethodNotAllowed)
		return
	}

    r.ParseForm();
    email := r.FormValue("email")
    password := r.FormValue("password")

    sessionid, err := login(db, email, password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        log.Println(err.Error())
        return
    }

    session := fmt.Sprint(sessionid)

    cookie := &http.Cookie{
		Name:    "bricked-up_login",
		Value:   session,
		Expires: time.Now().Add(12 * 30 * 24 * time.Hour),
		Secure:  true,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
}

// SignupHandler handles requests to the /signup endpoint.
// It restricts the request method to GET and responds with a placeholder message.
func SignupHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method unsupported", http.StatusMethodNotAllowed)
		return
	}

    r.ParseForm();
    email := r.FormValue("email")
    password := r.FormValue("password")

    err := registerUser(db, email, password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        log.Println(err.Error())
        return
    }
}

// VerifyHandler handles requests to the /verify endpoint.
// Only GET requests are supported, and it returns a placeholder response.
func VerifyHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method unsupported", http.StatusMethodNotAllowed)
		return
	}
	

	// verifyIDStr := r.URL.Query().Get("id")
	// code := r.URL.Query().Get("code")

	// if verifyIDStr == "" || code == "" {
	// 	http.Error(w, "Missing id or code", http.StatusBadRequest)
	// 	return
	// }

	// verifyID, err := strconv.Atoi(verifyIDStr)
	// if err != nil {
	// 	http.Error(w, "Invalid verification ID", http.StatusBadRequest)
	// 	return
	// }

	// // Check if the verify code exists and matches
	// var userID int
	// err = db.QueryRow(`SELECT userid FROM VERIFY WHERE id = ? AND code = ?`, verifyID, code).Scan(&userID)
	// if err == sql.ErrNoRows {
	// 	http.Error(w, "Invalid verification code or ID", http.StatusUnauthorized)
	// 	return
	// } else if err != nil {
	// 	http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// // Update user record to set verifyid
	// _, err = db.Exec(`UPDATE USER SET verifyid = ? WHERE id = ?`, verifyID, userID)
	// if err != nil {
	// 	http.Error(w, "Failed to update user verification: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// // Delete the verification entry
	// _, err = db.Exec(`DELETE FROM VERIFY WHERE id = ?`, verifyID)
	// if err != nil {
	// 	http.Error(w, "Failed to clean up verification record: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// // Respond with success
	// w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "TODO: verify")
}

// ChangeNameHandler updates the display name for the logged-in user.
func ChangeNameHandler (db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w,"Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()

	if err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("bricked-up_login")
	if err != nil {
		http.Error(w, "Invalid cookie session", http.StatusUnauthorized)
		return
	}

	sessionID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		http.Error(w,"Invalid session ID", http.StatusBadRequest)
		return
	}

	newName := r.FormValue("name")
	if newName == "" {
		http.Error(w, "Name cannot be empty", http.StatusBadRequest)
		return
	}

	err = ChangeDisplayName(db, sessionID, newName) 
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Name succesfully changed to %s", newName)
}

// CreateIssueHandler inserts a new issue into the database.
func CreateIssueHandler (db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w,"Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Extract and validate form values
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
	issueID, err := CreateIssue(title, desc, tagID, priority, completed, cost, date, db)
	if err != nil {
		http.Error(w, "Failed to create issue: "+err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Issue created successfully with ID %d", issueID)
}

// createTagHandler handles POST requests to create a new tag associated with a project.
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
	cookie, err := r.Cookie("bricked-up_login")
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
	tagID, err := CreateTag(db, sessionID, projectID, tagName, tagColor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("createTag error:", err)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Tag created successfully with ID %d", tagID)
}

// deleteUserHandler handles POST requests to delete the currently logged-in user.
// It reads the session cookie, validates it, and calls deleteUser to remove all user data.
func DeleteUserHandler (db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
        http.Error(w, "Invalid form data", http.StatusBadRequest)
        return
    }

	// Get session cookie
	cookie, err := r.Cookie("bricked-up_login")
	if err != nil {
		http.Error(w, "Missing session cookie", http.StatusUnauthorized)
		return
	}

	sessionID := cookie.Value

	// Call backend logic to delete the user
	if err := deleteUser(db, sessionID); err != nil {
		http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
		log.Println("deleteUser error:", err)
		return
	}

	// Invalidate the cookie
	cleared := &http.Cookie{
		Name:     "bricked-up_login",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // deletes the cookie
		HttpOnly: true,
	}
	http.SetCookie(w, cleared)

	// Respond with success
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "User successfully deleted")
}

// deleteTagHandler handles POST requests to delete a tag by its ID.
// The session is used to validate whether the user has permission to delete the tag.
func DeleteTagHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form values
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Retrieve session ID from cookie
	cookie, err := r.Cookie("bricked-up_login")
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
	err = DeleteTag(db, sessionID, tagID)
	if err != nil {
		http.Error(w, "Failed to delete tag: "+err.Error(), http.StatusForbidden)
		log.Println("deleteTag error:", err)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Tag with ID %d successfully deleted", tagID)
}

func GetIssueDetailsHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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
	jsonStr, err := getIssueDetails(db, issueID)
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

// getOrgMembersHandler handles GET requests to retrieve all user IDs 
// belonging to a specific organization.
// It expects an `orgid` query parameter, then returns a JSON array.
func GetOrgMembersHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse `orgid` from query parameters
	orgIDStr := r.URL.Query().Get("orgid")
	if orgIDStr == "" {
		http.Error(w, "Missing orgid parameter", http.StatusBadRequest)
		return
	}

	orgID, err := strconv.Atoi(orgIDStr)
	if err != nil {
		http.Error(w, "Invalid orgid", http.StatusBadRequest)
		return
	}

	// Call the core logic function
	jsonResult, err := GetOrgMembers(db, orgID)
	if err != nil {
		http.Error(w, "Failed to get org members: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, jsonResult)
}

// createOrganizationHandler handles POST requests to create a new organization
// and assigns the user (from the session) as an admin with full privileges.
func CreateOrganizationHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Parse form data
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Invalid form data", http.StatusBadRequest)
        return
    }

    // Retrieve session ID from cookie
    cookie, err := r.Cookie("bricked-up_login")
    if err != nil {
        http.Error(w, "Missing session cookie", http.StatusUnauthorized)
        return
    }

    sessionID, err := strconv.Atoi(cookie.Value)
    if err != nil {
        http.Error(w, "Invalid session ID", http.StatusBadRequest)
        return
    }

    // Get orgName from form
    orgName := r.FormValue("orgName")
    if orgName == "" {
        http.Error(w, "Missing orgName", http.StatusBadRequest)
        return
    }

    // Call backend logic
    orgID, err := CreateOrganization(db, sessionID, orgName)
    if err != nil {
        http.Error(w, "Failed to create organization: "+err.Error(), http.StatusInternalServerError)
        log.Println("createOrganization error:", err)
        return
    }

    // Return success with the new organization ID
    w.WriteHeader(http.StatusCreated)
    fmt.Fprintf(w, "Organization created successfully with ID %d", orgID)
}

// deleteOrganizationHandler handles POST requests to delete an organization.
// It requires the user to have admin (exec) privileges in the organization.
func DeleteOrganizationHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Parse form data
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Invalid form data", http.StatusBadRequest)
        return
    }

    // Retrieve session ID from cookie
    cookie, err := r.Cookie("bricked-up_login")
    if err != nil {
        http.Error(w, "Missing session cookie", http.StatusUnauthorized)
        return
    }

    sessionID, err := strconv.Atoi(cookie.Value)
    if err != nil {
        http.Error(w, "Invalid session ID", http.StatusBadRequest)
        return
    }

    // Parse organization ID
    orgIDStr := r.FormValue("orgid")
    if orgIDStr == "" {
        http.Error(w, "Missing orgid parameter", http.StatusBadRequest)
        return
    }

    orgID, err := strconv.Atoi(orgIDStr)
    if err != nil {
        http.Error(w, "Invalid orgid", http.StatusBadRequest)
        return
    }

    // Call the core DeleteOrganization logic
    err = DeleteOrganization(db, sessionID, orgID)
    if err != nil {
        // Check for known error types or just return internal server error
        log.Println("deleteOrganization error:", err)
        // You could differentiate error messages, e.g.:
        // - "organization does not exist" -> 404
        // - "user does not have permission" -> 403
        // Here we’ll do a generic 403 if it's a permission or existence issue:
        if err.Error() == "organization does not exist" ||
           err.Error() == "no session exists for the provided sessionID" ||
           err.Error() == "user is not a member of this organization" ||
           err.Error() == "user does not have permission to delete the organization" {
            http.Error(w, err.Error(), http.StatusForbidden)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Return success
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Organization with ID %d was deleted successfully", orgID)
}

// removeOrgMemberRoleHandler handles the removal of a role from a user in an organization.
// It expects a POST request with `orgMemberRoleId` in the form and a valid session cookie.
func RemoveOrgMemberRoleHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    // Only allow POST
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Parse form data
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Invalid form data", http.StatusBadRequest)
        return
    }

    // Retrieve session ID from cookie
    cookie, err := r.Cookie("bricked-up_login")
    if err != nil {
        http.Error(w, "Missing session cookie", http.StatusUnauthorized)
        return
    }

    sessionID, err := strconv.Atoi(cookie.Value)
    if err != nil {
        http.Error(w, "Invalid session ID", http.StatusBadRequest)
        return
    }

    // Parse orgMemberRoleId from the form
    roleIDStr := r.FormValue("orgMemberRoleId")
    if roleIDStr == "" {
        http.Error(w, "Missing orgMemberRoleId", http.StatusBadRequest)
        return
    }

    orgMemberRoleID, err := strconv.Atoi(roleIDStr)
    if err != nil {
        http.Error(w, "Invalid orgMemberRoleId", http.StatusBadRequest)
        return
    }

    // Call the core logic to remove the role
    err = RemoveOrgMemberRole(db, sessionID, orgMemberRoleID)
    if err != nil {
        // You might want more granular error handling here if needed:
        // e.g., 403 Forbidden for permission errors, 404 if role not found, etc.
        log.Println("removeOrgMemberRole error:", err)
        http.Error(w, "Failed to remove org member role: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Respond with success
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Successfully removed role assignment with ID %d", orgMemberRoleID)
}

// assignOrgRoleHandler handles POST requests to assign a new role (newRoleID)
// to User B (userID) in an organization (orgID).
// The acting user is determined by the session cookie (sessionID).
func AssignOrgRoleHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Parse form data
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Invalid form data", http.StatusBadRequest)
        return
    }

    // Retrieve session ID from cookie (User A)
    cookie, err := r.Cookie("bricked-up_login")
    if err != nil {
        http.Error(w, "Missing session cookie", http.StatusUnauthorized)
        return
    }
    sessionID, err := strconv.Atoi(cookie.Value)
    if err != nil {
        http.Error(w, "Invalid session ID", http.StatusBadRequest)
        return
    }

    // Parse form fields for userID, orgID, and newRoleID (User B, organization, role)
    userIDStr := r.FormValue("userID")
    if userIDStr == "" {
        http.Error(w, "Missing userID parameter", http.StatusBadRequest)
        return
    }
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid userID", http.StatusBadRequest)
        return
    }

    orgIDStr := r.FormValue("orgID")
    if orgIDStr == "" {
        http.Error(w, "Missing orgID parameter", http.StatusBadRequest)
        return
    }
    orgID, err := strconv.Atoi(orgIDStr)
    if err != nil {
        http.Error(w, "Invalid orgID", http.StatusBadRequest)
        return
    }

    newRoleIDStr := r.FormValue("newRoleID")
    if newRoleIDStr == "" {
        http.Error(w, "Missing newRoleID parameter", http.StatusBadRequest)
        return
    }
    newRoleID, err := strconv.Atoi(newRoleIDStr)
    if err != nil {
        http.Error(w, "Invalid newRoleID", http.StatusBadRequest)
        return
    }

    // Attempt to assign the role
    err = assignOrgRole(db, sessionID, userID, orgID, newRoleID)
    if err != nil {
        log.Println("assignOrgRole error:", err)
        // Here you could do more nuanced checks for permissions (403) vs. not found (404).
        // We'll simply respond with 403 for known permission-like issues, else 500.
        http.Error(w, err.Error(), http.StatusForbidden)
        return
    }

    // Success
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "User %d has been assigned to role %d in organization %d", userID, newRoleID, orgID)
}

// MainHandler checks if the request URL matches a known endpoint.
// If it does, the corresponding handler is called; otherwise, it returns a 404 error.
func MainHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if handler, ok := endpoints[r.URL.Path]; ok {
		handler(db, w, r)
		return
	}
	http.NotFound(w, r)
}
