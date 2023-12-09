package main

import (
	"github.com/vertex-center/vertex/apps/sql"
	"github.com/vertex-center/vertex/common/app"
)

func main() {
	app.RunStandalone(sql.NewApp())
}
