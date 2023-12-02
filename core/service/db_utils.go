package service

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

func (s *DbService) copyDb(from *sqlx.DB, to *sqlx.DB) error {
	tables := []string{
		"admin_settings",
		"migrations",
		"users",
		"credentials_argon2",
		"credentials_argon2_users",
		"sessions",
	}

	toTx, err := to.Beginx()
	if err != nil {
		return err
	}

	for _, t := range tables {
		err := s.copyDbTable(t, from, toTx)
		if err != nil {
			_ = toTx.Rollback()
			return err
		}
	}

	return toTx.Commit()
}

func (s *DbService) copyDbTable(name string, from *sqlx.DB, to *sqlx.Tx) error {
	log.Info("copying table",
		vlog.String("table", name),
		vlog.String("from", from.DriverName()),
		vlog.String("to", to.DriverName()),
	)

	rows, err := from.Queryx(fmt.Sprintf("SELECT * FROM %s", name))
	if err != nil {
		return err
	}
	defer rows.Close()

	// Delete all rows in the table
	_, err = to.Exec(fmt.Sprintf("DELETE FROM %s", name))
	if err != nil {
		return err
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	for rows.Next() {
		var values []string

		row := map[string]interface{}{}
		err := rows.MapScan(row)
		if err != nil {
			return err
		}

		for _, c := range columns {
			v := row[c]
			if v == nil {
				values = append(values, "NULL")
				continue
			}
			if _, ok := v.(string); ok {
				values = append(values, fmt.Sprintf("'%s'", v))
				continue
			}
			values = append(values, fmt.Sprintf("%v", v))
		}

		q := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			name,
			strings.Join(columns, ", "),
			strings.Join(values, ", "),
		)

		log.Debug(q)

		_, err = to.Exec(q)
		if err != nil {
			return err
		}
	}

	return nil
}
