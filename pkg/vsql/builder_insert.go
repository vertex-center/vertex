package vsql

import (
	"fmt"
	"strings"
)

type QueryInsert struct {
	table   string
	columns []string
	values  []string
}

func InsertInto(table string) *QueryInsert {
	return &QueryInsert{
		table: table,
	}
}

func (q *QueryInsert) Columns(columns ...string) *QueryInsert {
	q.columns = columns
	return q
}

func (q *QueryInsert) Values(values ...interface{}) *QueryInsert {
	for _, v := range values {
		if t, ok := v.(string); ok {
			v = fmt.Sprintf("'%s'", t)
		} else if t, ok := v.(bool); ok {
			v = fmt.Sprintf("%t", t)
		}
		q.values = append(q.values, fmt.Sprintf("%v", v))
	}
	return q
}

func (q *QueryInsert) Build(driver Driver) string {
	columns := strings.Join(q.columns, ", ")
	values := strings.Join(q.values, ", ")
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", q.table, columns, values)
}
