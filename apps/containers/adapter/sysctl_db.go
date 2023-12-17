package adapter

import (
	"context"
	"database/sql"

	"github.com/juju/errors"
	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/common/storage"
)

type sysctlDBAdapter struct {
	db storage.DB
}

func NewSysctlDBAdapter(db storage.DB) port.SysctlAdapter {
	return &sysctlDBAdapter{db}
}

func (a *sysctlDBAdapter) GetSysctls(ctx context.Context, id uuid.UUID) (types.Sysctls, error) {
	var sysctls types.Sysctls
	err := a.db.Select(&sysctls, `
		SELECT * FROM sysctls
		WHERE container_id = $1
	`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return sysctls, nil
	}
	return sysctls, err
}

func (a *sysctlDBAdapter) CreateSysctl(ctx context.Context, sysctl types.Sysctl) error {
	_, err := a.db.NamedExec(`
		INSERT INTO sysctls (container_id, name, value)
		VALUES (:container_id, :name, :value)
	`, sysctl)
	return err
}

func (a *sysctlDBAdapter) DeleteSysctls(ctx context.Context, id uuid.UUID) error {
	_, err := a.db.Exec(`
		DELETE FROM sysctls
		WHERE container_id = $1
	`, id)
	return err
}
