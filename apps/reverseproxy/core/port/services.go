package port

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/reverseproxy/core/types"
)

type (
	ProxyService interface {
		GetRedirects() types.ProxyRedirects
		GetRedirectByHost(host string) *types.ProxyRedirect
		AddRedirect(redirect types.ProxyRedirect) error
		RemoveRedirect(id uuid.UUID) error
	}
)
