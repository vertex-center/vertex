package types

type SqlDBMS struct {
	Username  string         `json:"username"`
	Password  string         `json:"password"`
	Databases *[]SqlDatabase `json:"databases,omitempty"`
}

type SqlDatabase struct {
	Name string `json:"name"`
}

type SqlDBMSAdapterPort interface {
	// GetDatabases returns a list of databases available in the DBMS.
	// If the DBMS does not support this operation, it will return nil.
	// If there is no database available, it will return an empty list.
	GetDatabases() (*[]SqlDatabase, error)
}
