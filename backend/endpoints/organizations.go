package endpoints

import (
	"brickedup/backend/organizations"
	"brickedup/backend/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// GetOrg handles GET requests to retrieve all information
// about a specific organization on /get-org.
// It expects an `orgid` URL query parameter, then returns a JSON array.
func GetOrgHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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
		log.Println(err.Error())
		return
	}

	// Call the core logic function
	org, err := organizations.GetOrg(db, orgID)
	if err != nil {
		http.Error(w, "Failed to get org: "+err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	json, err := json.Marshal(org)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// CreateOrganizationHandler handles POST requests to create a new organization
// on /create-org.
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

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	session := r.FormValue("sessionid")
	sessionID, err := strconv.Atoi(session)
	if err != nil {
		http.Error(w,"Invalid session ID", http.StatusBadRequest)
		return
	}

    // Get orgName from form
    orgName := r.FormValue("orgName")
    if orgName == "" {
        http.Error(w, "Missing orgName", http.StatusBadRequest)
        return
    }

    // Call backend logic
    _, err = organizations.CreateOrganization(db, sessionID, orgName)
    if err != nil {
        http.Error(w, "Failed to create organization: "+err.Error(), http.StatusInternalServerError)
        log.Println("createOrganization error:", err)
        return
    }

    // Return success with the new organization ID
    w.WriteHeader(http.StatusCreated)
}

// DeleteOrganizationHandler handles DELETE requests to delete an organization
// on /delete-org.
// It requires the user to have admin (exec) privileges in the organization.
func DeleteOrganizationHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Parse form data
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Invalid form data", http.StatusBadRequest)
        return
    }

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	session := r.FormValue("sessionid")
	sessionID, err := strconv.Atoi(session)
	if err != nil {
		http.Error(w,"Invalid session ID", http.StatusBadRequest)
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
    err = organizations.DeleteOrganization(db, sessionID, orgID)
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
}

// WithdrawOrgRoleHandler handles the withdrawal of a role from a user in an organization
// on /withdraw-org-role.
// It expects a DELETE request with `orgMemberRoleId` in the form and a valid session cookie.
func WithdrawOrgRoleHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Parse form data
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Invalid form data", http.StatusBadRequest)
        return
    }

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	session := r.FormValue("sessionid")
	sessionID, err := strconv.Atoi(session)
	if err != nil {
		http.Error(w,"Invalid session ID", http.StatusBadRequest)
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
    err = organizations.WithdrawOrgRole(db, sessionID, orgMemberRoleID)
    if err != nil {
        // You might want more granular error handling here if needed:
        // e.g., 403 Forbidden for permission errors, 404 if role not found, etc.
        log.Println("WithdrawOrgRole error:", err)
        http.Error(w, "Failed to withdraw org role: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Respond with success
    w.WriteHeader(http.StatusOK)
}

// AssignOrgRoleHandler handles POST requests to assign a new role (newRoleID)
// to User B (userID) in an organization (orgID) on /assign-org-role.
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

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	session := r.FormValue("sessionid")
	sessionID, err := strconv.Atoi(session)
	if err != nil {
		http.Error(w,"Invalid session ID", http.StatusBadRequest)
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
    err = organizations.AssignOrgRole(db, sessionID, userID, orgID, newRoleID)
    if err != nil {
        log.Println("assignOrgRole error:", err)
        // Here you could do more nuanced checks for permissions (403) vs. not found (404).
        // We'll simply respond with 403 for known permission-like issues, else 500.
        http.Error(w, err.Error(), http.StatusForbidden)
        return
    }

    // Success
    w.WriteHeader(http.StatusOK)
}

// GetOrgMemberHandler handles GET requests to retrieve information 
// about a organization member on /get-org-member.
// It takes `memberid` as a URL parameter.
func GetOrgMemberHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
        http.Error(w, "Invalid form data", http.StatusBadRequest)
        return
    }

	memberParam := r.URL.Query().Get("memberid")
	memberid, err := strconv.Atoi(memberParam)

	if err != nil {
        http.Error(w, "Invalid parameter for memberid", http.StatusBadRequest)
        return
	}


	user, err := organizations.GetOrgMember(db, memberid)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		log.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// GetOrgRole handles GET requests to retrieve information 
// about an organization role on /get-org-role.
// It takes `roleid` as a URL parameter.
func GetOrgRoleHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
        http.Error(w, "Invalid form data", http.StatusBadRequest)
        return
    }

	roleParam := r.URL.Query().Get("roleid")
	roleid, err := strconv.Atoi(roleParam)

	if err != nil {
        http.Error(w, "Invalid parameter for roleid", http.StatusBadRequest)
        return
	}


	role, err := organizations.GetOrgRole(db, roleid)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json, err := json.Marshal(role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		log.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// AddOrgMemberHandler handles POST requests to add a user to an organization on
// /add-org-member.
func AddOrgMemberHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
        http.Error(w, "Invalid form data", http.StatusBadRequest)
        return
    }

	session := r.FormValue("sessionid")
	sessionid, err := strconv.ParseInt(session, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	user := r.FormValue("userid")
	userid, err := strconv.Atoi(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	role := r.FormValue("roleid")
	roleid, err := strconv.Atoi(role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	org := r.FormValue("orgid")
	orgid, err := strconv.Atoi(org)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	err = organizations.AddOrgMember(db, sessionid, userid, roleid, orgid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

// RemoveOrgMemberHandler handles DELETE requests to remove a member from an organization
// on /remove-proj-member.
func RemoveOrgMemberHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
        http.Error(w, "Invalid form data", http.StatusBadRequest)
        return
    }

	session := r.FormValue("sessionid")
	sessionid, err := strconv.ParseInt(session, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	member := r.FormValue("memberid")
	memberid, err := strconv.Atoi(member)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	err = organizations.RemoveOrgMember(db, sessionid, memberid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateOrgHandler handles PATCH requests to update an organization on
// /update-org.
func UpdateOrgHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	session := r.FormValue("sessionid")
	sessionid, err := strconv.Atoi(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	org := r.FormValue("orgid")
	orgid, err := strconv.Atoi(org)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	name := r.FormValue("nameid")

	updated_org := utils.Organization {
		ID: orgid,
		Name: name,
	}

	err = organizations.UpdateOrg(db, sessionid, orgid, updated_org)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
