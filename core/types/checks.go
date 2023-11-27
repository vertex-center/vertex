package types

type CheckResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Error error  `json:"error"`
}
