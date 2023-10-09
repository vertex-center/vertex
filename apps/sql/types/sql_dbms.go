package types

type DBMS struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Databases *[]DB  `json:"databases,omitempty"`
}

type DBMSAdapterPort interface {
	// GetDatabases returns a list of databases available in the DBMS.
	// If the DBMS does not support this operation, it will return nil.
	// If there is no database available, it will return an empty list.
	GetDatabases() (*[]DB, error)
}
