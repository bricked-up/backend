package backend

// Project contains the details of a project.
type Project struct {
	ID       int		`json:"id"`
	OrgID    int		`json:"orgid"`
	Name     string		`json:"name"`
	Budget   int		`json:"budget"`
	Charter  string		`json:"charter"`
	Archived bool		`json:"archived"`
}

// Organization contains the details of an organization.
type Organization struct {
	ID    int 			`json:"id"`
	Name  string		`json:"name"`
}

