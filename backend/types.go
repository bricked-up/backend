package backend

// Organization represents the structure of the ORGANIZATION table in the database
type Org struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
