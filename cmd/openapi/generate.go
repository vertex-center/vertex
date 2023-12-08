package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/carlmjohnson/requests"
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/admin"
	"github.com/vertex-center/vertex/apps/auth"
	"github.com/vertex-center/vertex/apps/containers"
	"github.com/vertex-center/vertex/apps/monitoring"
	"github.com/vertex-center/vertex/apps/reverseproxy"
	"github.com/vertex-center/vertex/apps/serviceeditor"
	"github.com/vertex-center/vertex/apps/sql"
	"github.com/vertex-center/vertex/apps/tunnels"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/core/types/server"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
)

func main() {
	apps := []app.Interface{
		admin.NewApp(),
		auth.NewApp(),
		containers.NewApp(),
		monitoring.NewApp(),
		reverseproxy.NewApp(),
		serviceeditor.NewApp(),
		sql.NewApp(),
		tunnels.NewApp(),
	}

	err := os.MkdirAll("openapi", 0755)
	if err != nil {
		panic(err)
	}

	port := os.Getenv("VERTEX_OPENAPI_PORT")
	if port == "" {
		port = "6955"
	}

	u, err := url.Parse(fmt.Sprintf("http://localhost:%s", port))
	if err != nil {
		panic(err)
	}

	for _, a := range apps {
		if a, ok := a.(app.InitializableRouter); ok {
			srv := runServer(a.InitializeRouter, a.Meta(), u)
			downloadOpenAPI(a.Meta().ID, u)
			srv.Stop()
		}

		if a, ok := a.(app.KernelInitializableRouter); ok {
			srv := runServer(a.InitializeKernelRouter, a.Meta(), u)
			downloadOpenAPI(a.Meta().ID+"_kernel", u)
			srv.Stop()
		}
	}
}

func runServer(initRoutes func(r *fizz.RouterGroup) error, meta app.Meta, u *url.URL) *server.Server {
	vertexCtx := types.NewVertexContext(types.About{}, true)
	ctx := app.NewContext(vertexCtx)

	info := openapi.Info{
		Title:       meta.Name,
		Description: meta.Description,
		Version:     ctx.About().Version,
	}

	srv := server.New(meta.ID, &info, u, vertexCtx)

	base := srv.Router.Group("/api", "", "")
	err := initRoutes(base)
	if err != nil {
		panic(err)
	}

	base.GET("/ping", []fizz.OperationOption{
		fizz.Summary("Ping the app"),
	}, router.Handler(func(c *gin.Context) error {
		return nil
	}))

	_ = srv.StartAsync()

	return srv
}

func downloadOpenAPI(id string, u *url.URL) {
	w, err := os.Create(path.Join("openapi", "openapi."+id+".yaml"))
	if err != nil {
		panic(err)
	}
	defer w.Close()

	err = requests.New().
		BaseURL(u.String()).
		Path("/openapi.yaml").
		ToWriter(w).
		Fetch(context.Background())
	if err != nil {
		panic(err)
	}
}
