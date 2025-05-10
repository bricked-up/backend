package endpoints

import (
	"database/sql"
	"net/http"
)

// DBHandlerFunc is the function prototype for endpoint handlers.
type DBHandlerFunc func(*sql.DB, http.ResponseWriter, *http.Request)

// Endpoints maps URL paths to their corresponding handler functions.
var Endpoints = map[string]DBHandlerFunc{
	"/login":                   	LoginHandler,
	"/signup":                  	SignupHandler,
	"/verify":                  	VerifyHandler,
	"/get-user":               		GetUserHandler,
	"/get-all-users":          		GetAllUsersHandler,
	"/delete-user":            		DeleteUserHandler,
	"/update-user":            		UpdateUserHandler,
	"/create-issue":           		CreateIssueHandler,
	"/get-issue":               	GetIssueHandler,
	"/create-tag":             		CreateTagHandler,
	"/delete-tag":             		DeleteTagHandler,
	"/get-org":         			GetOrgHandler,
	"/get-org-member":    			GetOrgMemberHandler,
	"/get-org-role":				GetOrgRoleHandler,
	"/add-org-member":				AddOrgMemberHandler,
	"/create-org":             		CreateOrganizationHandler,
	"/delete-org":             		DeleteOrganizationHandler,
	"/withdraw-org-role":		 	WithdrawOrgRoleHandler,
	"/assign-org-role":        		AssignOrgRoleHandler,
	"/get-proj":					GetProjHandler,
	"/get-proj-member":				GetProjMemberHandler,
	"/get-proj-role":				GetProjRoleHandler,
	"/add-proj-member":				AddProjMemberHandler,
	"/remove-proj-member":			RemoveProjMemberHandler,
	"/get-tag":						GetTagHandler,
	"/archive-proj": 				ArchiveProjHandler,
}
