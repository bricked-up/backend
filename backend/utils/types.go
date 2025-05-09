package utils

import (
	"database/sql"
	"time"
)

// User contains the details of a user along with all projects
// and organizations that they are a part of.
type User struct {
	ID				int 		`json:"id"`
	Name 			string		`json:"name"`
	Email   		string 		`json:"email"`
	Password 		string 		`json:"password"`
	Avatar 			string 		`json:"avatar"`
	Verified    	bool		`json:"verified"`
	Projects		[]int 		`json:"projects"`
	Organizations	[]int		`json:"organizations"`
	Issues			[]int 		`json:"issues"`
}

// Project contains the details of a project.
type Project struct {
	ID       int		`json:"id"`
	OrgID    int		`json:"orgid"`
	Name     string		`json:"name"`
	Budget   int		`json:"budget"`
	Charter  string		`json:"charter"`
	Archived bool		`json:"archived"`
	Members []int 		`json:"members"`
	Issues  []int		`json:"issues"`
	Tags	[]int		`json:"tags"`
	Roles 	[]int		`json:"roles"`
}

// ProjectMember contains all information relating to the user in a given 
// project.
type ProjectMember struct {
	ID 			int			`json:"id"`
	UserID		int 		`json:"userid"`
	ProjectID 	int			`json:"projectid"`
	Roles 		[]int		`json:"roles"`
	CanExec		bool		`json:"can_exec"`
	CanWrite	bool		`json:"can_write"`
	CanRead		bool		`json:"can_read"`
	Issues 		[]int		`json:"issues"`
}

// ProjectRole contains all information relating to a role in a project.
type ProjectRole struct {
	ID 			int			`json:"id"`
	ProjectID	int 		`json:"projectid"`
	Name 		string 		`json:"name"`
	CanExec		bool		`json:"can_exec"`
	CanWrite	bool		`json:"can_write"`
	CanRead		bool		`json:"can_read"`
}

// Issue contains all information relating to an issue.
type Issue struct {
	ID       		int				`json:"id"`
	Title    		string			`json:"title"`
	Desc     		string			`json:"desc"`
	Cost			int				`json:"cost"`
	TagID    		int				`json:"tagid"`
	Priority 		int				`json:"priority"`
	Created  		time.Time		`json:"created"`
	Completed  		sql.NullTime	`json:"completed"`
	Dependencies	[]int			`json:"dependencies"`
}

// Organization contains the details of an organization.
type Organization struct {
	ID    		int 		`json:"id"`
	Name  		string		`json:"name"`
	Members 	[]int 		`json:"members"`
	Projects 	[]int 		`json:"projects"`
	Roles 		[]int		`json:"roles"`
}

// OrgMember contains all information relating to the user in a given 
// organization.
type OrgMember struct {
	ID 				int			`json:"id"`
	UserID			int 		`json:"userid"`
	OrganizationID 	int			`json:"projectid"`
	Roles 			[]int		`json:"roles"`
	CanExec			bool		`json:"can_exec"`
	CanWrite		bool		`json:"can_write"`
	CanRead			bool		`json:"can_read"`
}
