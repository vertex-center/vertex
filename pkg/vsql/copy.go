package vsql

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

// CopyTables copies the tables from one database to another.
// The tables are copied in one unique transaction.
func CopyTables(from *sqlx.DB, to *sqlx.DB, tables []string) error {
	toTx, err := to.Beginx()
	if err != nil {
		return err
	}

	for _, t := range tables {
		err := CopyTable(from, toTx, t)
		if err != nil {
			_ = toTx.Rollback()
			return err
		}
	}

	return toTx.Commit()
}

// CopyTable copies the given table from one database to another.
func CopyTable(from *sqlx.DB, to *sqlx.Tx, name string) error {
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
