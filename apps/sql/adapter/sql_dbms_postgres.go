package adapter

import (
	"errors"
	"github.com/vertex-center/vertex/apps/sql/core/port"
	types2 "github.com/vertex-center/vertex/apps/sql/core/types"
	"strconv"

	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SqlDBMSPostgresAdapter struct {
	*SqlDBMSAdapter

	db *gorm.DB
}

type SqlDBMSPostgresAdapterParams struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewSqlDBMSPostgresAdapter(params *SqlDBMSPostgresAdapterParams) port.DBMSAdapter {
	adapter := &SqlDBMSPostgresAdapter{
		SqlDBMSAdapter: NewSqlDBMSAdapter().(*SqlDBMSAdapter),
	}

	go func() {
		dns := "host=localhost"
		if params.Port != 0 {
			dns += " port=" + strconv.Itoa(params.Port)
		}
		if params.Username != "" {
			dns += " user=" + params.Username
		}
		if params.Password != "" {
			dns += " password=" + params.Password
		}
		dns += " dbname=postgres sslmode=disable"

		log.Info("connecting to postgres", vlog.String("dns", dns))

		var err error
		adapter.db, err = gorm.Open(postgres.Open(dns), &gorm.Config{
			DisableAutomaticPing: true,
		})
		if err != nil {
			log.Error(err)
		}
	}()

	return adapter
}

func (a *SqlDBMSPostgresAdapter) GetDatabases() (*[]types2.DB, error) {
	if a.db == nil {
		return nil, errors.New("connection not established")
	}

	rows, err := a.db.Table("pg_database").Where("datistemplate = false").Select("datname").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []types2.DB
	for rows.Next() {
		var database types2.DB
		err := rows.Scan(&database.Name)
		if err != nil {
			return nil, err
		}

		databases = append(databases, database)
	}
	return &databases, nil
}
