package vsql

import (
	"fmt"
	"strings"
)

type QueryCreateTable struct {
	name   string
	fields []Builder
}

func CreateTable(name string) *QueryCreateTable {
	return &QueryCreateTable{
		name: name,
	}
}

func (q *QueryCreateTable) WithID() *QueryCreateTable {
	q.fields = append(q.fields, &QueryWithID{})
	return q
}

func (q *QueryCreateTable) WithField(name string, dataType string, options ...string) *QueryCreateTable {
	q.fields = append(q.fields, &QueryWithField{
		name:     name,
		dataType: dataType,
		options:  options,
	})
	return q
}

func (q *QueryCreateTable) WithCreatedAt() *QueryCreateTable {
	q.fields = append(q.fields, &QueryWithCreatedAt{})
	return q
}

func (q *QueryCreateTable) WithUpdatedAt() *QueryCreateTable {
	q.fields = append(q.fields, &QueryWithUpdatedAt{})
	return q
}

func (q *QueryCreateTable) WithDeletedAt() *QueryCreateTable {
	q.fields = append(q.fields, &QueryWithDeletedAt{})
	return q
}

func (q *QueryCreateTable) WithPrimaryKey(fields ...string) *QueryCreateTable {
	q.fields = append(q.fields, &QueryWithPrimaryKey{
		fields: fields,
	})
	return q
}

func (q *QueryCreateTable) WithForeignKey(field string, table string, reference string) *QueryCreateTable {
	q.fields = append(q.fields, &QueryWithForeignKey{
		field:     field,
		table:     table,
		reference: reference,
	})
	return q
}

func (q *QueryCreateTable) Build(driver Driver) string {
	res := fmt.Sprintf("CREATE TABLE %s (", q.name)
	var fields []string
	for _, f := range q.fields {
		fields = append(fields, f.Build(driver))
	}
	res += strings.Join(fields, ", ")
	res += ");"
	return res
}

// Fields

type QueryWithID struct{}

func (q *QueryWithID) Build(driver Driver) string {
	return fmt.Sprintf("id %s", driver.AutoIncrement())
}

type QueryWithField struct {
	name     string
	dataType string
	options  []string
}

func (q *QueryWithField) Build(driver Driver) string {
	return fmt.Sprintf("%s %s %s", q.name, q.dataType, strings.Join(q.options, " "))
}

type QueryWithCreatedAt struct{}

func (q *QueryWithCreatedAt) Build(driver Driver) string {
	return "created_at INTEGER NOT NULL"
}

type QueryWithUpdatedAt struct{}

func (q *QueryWithUpdatedAt) Build(driver Driver) string {
	return "updated_at INTEGER NOT NULL"
}

type QueryWithDeletedAt struct{}

func (q *QueryWithDeletedAt) Build(driver Driver) string {
	return "deleted_at INTEGER DEFAULT NULL"
}

type QueryWithPrimaryKey struct {
	fields []string
}

func (q *QueryWithPrimaryKey) Build(driver Driver) string {
	return fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(q.fields, ", "))
}

type QueryWithForeignKey struct {
	field     string
	table     string
	reference string
}

func (q *QueryWithForeignKey) Build(driver Driver) string {
	return fmt.Sprintf(driver.ForeignKey(),
		q.field,
		q.table,
		q.reference,
	)
}
