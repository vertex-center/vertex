package main

import (
	"github.com/vertex-center/vertex/apps/sql"
	"github.com/vertex-center/vertex/core/types/app"
)

func main() {
	app.RunStandalone(sql.NewApp())
}
