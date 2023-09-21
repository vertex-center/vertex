package types

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return e.Message
}
