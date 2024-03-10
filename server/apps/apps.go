package apps

import (
	"github.com/vertex-center/vertex/server/apps/admin"
	"github.com/vertex-center/vertex/server/apps/auth"
	"github.com/vertex-center/vertex/server/apps/containers"
	"github.com/vertex-center/vertex/server/apps/logs"
	"github.com/vertex-center/vertex/server/apps/monitoring"
	"github.com/vertex-center/vertex/server/apps/reverseproxy"
	"github.com/vertex-center/vertex/server/apps/sql"
	"github.com/vertex-center/vertex/server/apps/tunnels"
	"github.com/vertex-center/vertex/server/common/app"
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
