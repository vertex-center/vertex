package types

import (
	"errors"

	"github.com/vertex-center/vertex/common/baseline"
)

var (
	ErrAlreadyUpdating = errors.New("an update is already in progress, cannot start another")
)

type Update struct {
	Baseline baseline.Baseline `json:"baseline"` // Baseline is the set of versions that are available to update to.
	Updating bool              `json:"updating"` // Updating is true if an update is currently in progress.
}
