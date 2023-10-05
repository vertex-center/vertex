package router

type ErrCode string

const (
	ErrFailedToParseBody ErrCode = "failed_to_parse_body"
)

type Error struct {
	Code           ErrCode `json:"code"`
	PublicMessage  string  `json:"message,omitempty"`
	PrivateMessage string  `json:"-"`
}

func (e Error) Error() string {
	return e.PublicMessage + "; " + e.PrivateMessage
}
