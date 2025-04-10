package endpoints

import (
	"database/sql"
	"net/http"
)

// DBHandlerFunc is the function prototype for endpoint handlers.
type DBHandlerFunc func(*sql.DB, http.ResponseWriter, *http.Request)

// LoginCookie specifies the name of the cookie that holds the user's session
// token.
const LoginCookie = "bricked-up_login"

// Endpoints maps URL paths to their corresponding handler functions.
var Endpoints = map[string]DBHandlerFunc{
	"/login":                   LoginHandler,
	"/signup":                  SignupHandler,
	"/verify":                  VerifyHandler,
	"/get-user":               	GetUserHandler,
	"/delete-user":            	DeleteUserHandler,
	"/update-user":            	UpdateUserHandler,
	"/create-issue":           	CreateIssueHandler,
	"/get-issue":               GetIssueHandler,
	"/create-tag":             	CreateTagHandler,
	"/delete-tag":             	DeleteTagHandler,
	"/get-org-members":         GetOrgMembersHandler,
	"/create-org":             	CreateOrganizationHandler,
	"/delete-org":             	DeleteOrganizationHandler,
	"/remove-org-member-role": 	RemoveOrgMemberRoleHandler,
	"/assign-org-role":        	AssignOrgRoleHandler,
}
