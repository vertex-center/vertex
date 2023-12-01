package vsql

func BuildSchema(driver Driver, queries ...Builder) string {
	schema := ""
	for _, q := range queries {
		schema += q.Build(driver)
	}
	return schema
}
