package utils

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

// Organization contains the details of an organization.
type Organization struct {
	ID    		int 		`json:"id"`
	Name  		string		`json:"name"`
	Members 	[]int 		`json:"members"`
	Projects 	[]int 		`json:"projects"`
	Roles 		[]int		`json:"roles"`
}

