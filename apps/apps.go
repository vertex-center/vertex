package apps

import (
	"github.com/vertex-center/vertex/apps/admin"
	"github.com/vertex-center/vertex/apps/auth"
	"github.com/vertex-center/vertex/apps/containers"
	"github.com/vertex-center/vertex/apps/logs"
	"github.com/vertex-center/vertex/apps/monitoring"
	"github.com/vertex-center/vertex/apps/reverseproxy"
	"github.com/vertex-center/vertex/apps/sql"
	"github.com/vertex-center/vertex/apps/tunnels"
	"github.com/vertex-center/vertex/common/app"
)

var Apps = []app.Interface{
	admin.NewApp(),
	auth.NewApp(),
	sql.NewApp(),
	tunnels.NewApp(),
	logs.NewApp(),
	monitoring.NewApp(),
	containers.NewApp(),
	reverseproxy.NewApp(),
}
