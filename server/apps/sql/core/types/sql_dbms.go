package types

type DBMS struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Databases *[]DB  `json:"databases,omitempty"`
}
