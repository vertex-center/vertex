package types

import "github.com/google/uuid"

type ProxyRedirects map[uuid.UUID]ProxyRedirect

type ProxyRedirect struct {
	Source string `json:"source" validate:"required"`
	Target string `json:"target" validate:"required"`
}
