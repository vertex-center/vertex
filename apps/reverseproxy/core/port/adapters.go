package port

import (
	"github.com/vertex-center/vertex/apps/reverseproxy/core/types"
	"github.com/vertex-center/vertex/common/uuid"
)

type ProxyAdapter interface {
	GetRedirects() types.ProxyRedirects
	GetRedirectByHost(host string) *types.ProxyRedirect
	AddRedirect(id uuid.UUID, redirect types.ProxyRedirect) error
	RemoveRedirect(id uuid.UUID) error
}
