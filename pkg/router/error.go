package router

import "errors"

type ErrCode string

var (
	ErrFailedToStopServer = errors.New("failed to stop the server")
)

type Error struct {
	Code           ErrCode `json:"code"`
	PublicMessage  string  `json:"message,omitempty"`
	PrivateMessage string  `json:"-"`
}

func (e Error) Error() string {
	return e.PublicMessage + "; " + e.PrivateMessage
}
