package port

import (
	"context"

	"github.com/vertex-center/vertex/server/apps/containers/core/types"
	sqltypes "github.com/vertex-center/vertex/server/apps/sql/core/types"
)

type SqlService interface {
	Get(inst *types.Container) (sqltypes.DBMS, error)
	Install(ctx context.Context, dbms string) (types.Container, error)
}
