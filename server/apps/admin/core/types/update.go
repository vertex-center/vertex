package types

import (
	"github.com/vertex-center/vertex/server/common/baseline"
)

type Update struct {
	Baseline baseline.Baseline `json:"baseline"` // Baseline is the set of versions that are available to update to.
}
