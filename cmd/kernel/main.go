package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/adapter"
	"github.com/vertex-center/vertex/apps/containers"
	"github.com/vertex-center/vertex/apps/monitoring"
	"github.com/vertex-center/vertex/apps/reverseproxy"
	"github.com/vertex-center/vertex/apps/serviceeditor"
	"github.com/vertex-center/vertex/apps/sql"
	"github.com/vertex-center/vertex/apps/tunnels"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/service"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/handler"
	"github.com/vertex-center/vertex/pkg/ginutils"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/netcap"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vlog"
)

var (
	r   *router.Router
	ctx *types.VertexContext

	sshAdapter port.SshAdapter

	sshService port.SshService
)

func main() {
	ensureRoot()
	parseArgs()

	// If go.mod is there, build vertex first.
	_, err := os.Stat("go.mod")
	if err == nil {
		log.Info("init.go found. Building vertex...")
		buildVertex()
	}

	err = netcap.AllowPortsManagement("vertex")
	if err != nil {
		log.Error(err)
	}

	initRouter()
	initAdapters()
	initServices()
	initRoutes()

	ctx.DispatchEvent(types.EventServerStart{})

	shutdownChan := make(chan os.Signal, 1)
	go func() {
		startRouter()
		shutdownChan <- syscall.SIGINT
	}()

	// Vertex
	var vertex *exec.Cmd
	go func() {
		var err error
		vertex, err = runVertex([]string{
			"-host", config.KernelCurrent.Host,
			"-port", config.KernelCurrent.Port,
			"-port-kernel", config.KernelCurrent.PortKernel,
			"-port-proxy", config.KernelCurrent.PortProxy,
			"-port-prometheus", config.KernelCurrent.PortPrometheus,
		}...)
		if err != nil {
			log.Error(err)
			shutdownChan <- syscall.SIGINT
			return
		}
		err = vertex.Wait()
		if err != nil {
			log.Error(err)
		}
		shutdownChan <- syscall.SIGINT
	}()

	// OS interrupt
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	<-shutdownChan
	log.Info("Shutting down...")
	if vertex != nil && vertex.Process != nil {
		_ = vertex.Process.Signal(os.Interrupt)
		_, _ = vertex.Process.Wait()
	}
	stopRouter()
}

func ensureRoot() {
	if os.Getuid() != 0 {
		log.Warn("vertex-kernel must be run as root to work properly")
	}
}

func parseArgs() {
	var (
		flagUsername       = flag.String("user", "", "username of the unprivileged user")
		flagUID            = flag.Uint("uid", 0, "uid of the unprivileged user")
		flagGID            = flag.Uint("gid", 0, "gid of the unprivileged user")
		flagHost           = flag.String("host", config.Current.Host, "The Vertex access url")
		flagPort           = flag.String("port", config.Current.Port, "The Vertex port")
		flagPortKernel     = flag.String("port-kernel", config.Current.PortKernel, "The Vertex Kernel port")
		flagPortProxy      = flag.String("port-proxy", config.Current.PortProxy, "The Vertex Proxy port")
		flagPortPrometheus = flag.String("port-prometheus", config.Current.PortPrometheus, "The Prometheus port")
	)

	flag.Parse()

	config.KernelCurrent.Host = *flagHost
	config.KernelCurrent.Port = *flagPort
	config.KernelCurrent.PortKernel = *flagPortKernel
	config.KernelCurrent.PortProxy = *flagPortProxy
	config.KernelCurrent.PortPrometheus = *flagPortPrometheus

	if *flagUsername == "" {
		*flagUsername = os.Getenv("USER")
		log.Warn("no username specified; trying to retrieve username from env", vlog.String("user", *flagUsername))
	}

	if *flagUsername != "" {
		u, err := user.Lookup(*flagUsername)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		uid, err := strconv.ParseInt(u.Uid, 10, 32)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		gid, err := strconv.ParseInt(u.Gid, 10, 32)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		config.KernelCurrent.Uid = uint32(uid)
		config.KernelCurrent.Gid = uint32(gid)
		return
	}

	config.KernelCurrent.Uid = uint32(*flagUID)
	config.KernelCurrent.Gid = uint32(*flagGID)
}

func buildVertex() {
	log.Info("Building vertex")

	start := time.Now()
	cmd := exec.Command("go", "build", "-o", "vertex", "cmd/main/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	end := time.Now()

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Info("Build completed in " + end.Sub(start).String())
}

func initRouter() {
	gin.SetMode(gin.ReleaseMode)
	ctx = types.NewVertexContext()
	r = router.New()
	r.Use(ginutils.ErrorHandler())
	r.Use(ginutils.Logger("KERNEL"))
	r.Use(gin.Recovery())
}

func initAdapters() {
	sshAdapter = adapter.NewSshFsAdapter(nil)
}

func initServices() {
	sshService = service.NewSshKernelService(sshAdapter)

	service.NewAppsService(ctx, true, r,
		[]app.Interface{
			sql.NewApp(),
			tunnels.NewApp(),
			monitoring.NewApp(),
			containers.NewApp(),
			reverseproxy.NewApp(),
			serviceeditor.NewApp(),
		},
	)
}

func initRoutes() {
	// docapi:k title Vertex Kernel
	// docapi:k description A platform to manage your self-hosted server.
	// docapi:k version 0.0.0
	// docapi:k filename kernel

	// docapi:k url http://{ip}:{port-kernel}/api
	// docapi:k urlvar ip localhost The IP address of the kernel.
	// docapi:k urlvar port-kernel 6131 The port of the server.

	api := r.Group("/api")

	sshHandler := handler.NewSshKernelHandler(sshService)
	ssh := api.Group("/security/ssh")
	// docapi:k route /security/ssh get_ssh_keys_kernel
	ssh.GET("", sshHandler.Get)
	// docapi:k route /security/ssh add_ssh_key_kernel
	ssh.POST("", sshHandler.Add)
	// docapi:k route /security/ssh/{fingerprint} delete_ssh_key_kernel
	ssh.DELETE("/:fingerprint", sshHandler.Delete)
}

func startRouter() {
	log.Info("vertex-kernel started", vlog.String("url", config.KernelCurrent.KernelURL()))
	addr := fmt.Sprintf(":%s", config.KernelCurrent.PortKernel)

	err := r.Start(addr)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func stopRouter() {
	ctx.DispatchEvent(types.EventServerStop{})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := r.Stop(ctx)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
