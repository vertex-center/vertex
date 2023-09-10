package types

import "github.com/google/uuid"

type ProxyRedirects map[uuid.UUID]ProxyRedirect

type ProxyRedirect struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type ProxyAdapterPort interface {
	GetRedirects() ProxyRedirects
	AddRedirect(id uuid.UUID, redirect ProxyRedirect) error
	RemoveRedirect(id uuid.UUID) error
}
