package endpoints

import (
	"brickedup/backend/projects"
	"brickedup/backend/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// GetProjHandler handles GET requests to retrieve data about a project
// on /get-proj.
// It takes `projecid` as a URL parameter.
func GetProjHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
        http.Error(w, "Invalid form data", http.StatusBadRequest)
        return
    }

	projParam := r.URL.Query().Get("projectid")
	projectid, err := strconv.Atoi(projParam)

	if err != nil {
        http.Error(w, "Invalid parameter for projectid", http.StatusBadRequest)
        return
	}


	project, err := projects.GetProject(db, projectid)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(project)
}
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

// ArchiveProjHandler handles POST requests to archive a project by its ID on 
// /archive-proj.
// The session is used to validate whether the user has the necessary permissions.
func ArchiveProjHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	// Parse form values
	projectID, err := strconv.Atoi(r.FormValue("projectid"))
	if err != nil {
		http.Error(w, "Invalid or missing project ID", http.StatusBadRequest)
		return
	}

	err = projects.ArchiveProj(db, sessionID, projectID)
	if err != nil {
		http.Error(w, "Failed to archive project: "+err.Error(), http.StatusForbidden)
		log.Println("ArchiveProj error:", err)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
}

// GetProjMemberHandler handles GET requests to retrieve information 
// about a project member on /get-proj-member.
// It takes `memberid` as a URL parameter.
func GetProjMemberHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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
        http.Error(w, "Invalid parameter for projectid", http.StatusBadRequest)
        return
	}


	user, err := projects.GetProjMember(db, memberid)

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

// GetTagHandler handles GET requests to retrieve information about
// a tag on /get-tag.
// It takes `tagid` as a URL parameter.
func GetTagHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
        http.Error(w, "Invalid form data", http.StatusBadRequest)
        return
    }

	tagParam := r.URL.Query().Get("tagid")
	tagid, err := strconv.Atoi(tagParam)

	if err != nil {
        http.Error(w, "Invalid parameter for tagid", http.StatusBadRequest)
        return
	}

	tag, err := projects.GetTag(db, tagid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		log.Println(err.Error())
		return
	}

	json, err := json.Marshal(tag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// GetProjRoleHandler handles GET requests to retrieve information 
// about a project role on /get-proj-role.
// It takes `roleid` as a URL parameter.
func GetProjRoleHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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


	role, err := projects.GetProjRole(db, roleid)

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

// AddProjMemberHandler handles POST requests to add a user to a project on
// /add-proj-member.
func AddProjMemberHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

	project := r.FormValue("projectid")
	projectid, err := strconv.Atoi(project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	err = projects.AddProjMember(db, sessionid, userid, roleid, projectid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

// RemoveProjMemberHandler handles DELETE requests to remove a member from a project on
// /remove-proj-member.
func RemoveProjMemberHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

	err = projects.RemoveProjMember(db, sessionid, memberid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetAllProjHandler handles GET requests and returns all projectIDs as JSON
// on /get-all-project.
func GetAllProjHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	projids, err := projects.GetAllProj(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	json, err := json.Marshal(projids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// CreateProj handles POST requests to create a project on /create-proj.
func CreateProjHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
        http.Error(w, "Invalid form data", http.StatusBadRequest)
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

	name := r.FormValue("name")
	charter := r.FormValue("charter")

	budgetstr := r.FormValue("budget")
	budget, err := strconv.Atoi(budgetstr)

	if err != nil {
        http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	err = projects.CreateProj(db, sessionid, orgid, name, budget, charter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateProjHandler handles PATCH requests to update an project on
// /update-proj.
func UpdateProjHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

	proj := r.FormValue("projectid")
	projid, err := strconv.Atoi(proj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	name := r.FormValue("name")
	charter := r.FormValue("charter")
	budgetstr := r.FormValue("budget")

	budget, err := strconv.Atoi(budgetstr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	updated_org := utils.Project {
		ID: projid,
		Name: name,
		Budget: budget,
		Charter: charter,
	}

	err = projects.UpdateProject(db, sessionid, projid, updated_org)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

