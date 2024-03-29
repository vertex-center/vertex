package service

import (
	"context"
	"errors"

	"github.com/vertex-center/vertex/server/apps/admin/core/port"
	"github.com/vertex-center/vertex/server/apps/admin/core/types"
	"github.com/vertex-center/vertex/server/common/log"
	"github.com/vertex-center/vertex/server/pkg/net"
	"github.com/vertex-center/vlog"
	"go.uber.org/atomic"
)

// checksService is a service used to check if
// Vertex is ready to serve or not.
type checksService struct{}

func NewChecksService() port.ChecksService {
	return &checksService{}
}

// CheckAll checks if Vertex is ready to serve, by doing
// some checks like internet connection, api readiness, etc.
// It returns a channel of CheckResponse, which contains
// the result of each check. The channel is closed when all
// checks are done.
func (s *checksService) CheckAll(ctx context.Context) <-chan types.CheckResponse {
	checks := []func(ctx context.Context) types.CheckResponse{
		s.checkInternet,
	}

	resChan := make(chan types.CheckResponse, len(checks))

	remaining := atomic.NewInt32(int32(len(checks)))
	for _, check := range checks {
		check := check
		go func() {
			res := check(ctx)
			resChan <- res
			if res.Error != "" {
				log.Error(errors.New("component check: failed"), vlog.String("id", res.ID), vlog.String("name", res.Name), vlog.String("reason", res.Error))
			} else {
				log.Info("component check: ok", vlog.String("id", res.ID), vlog.String("name", res.Name))
			}
			if remaining.Dec() == 0 {
				log.Info("all components are checked")
				close(resChan)
			}
		}()
	}

	return resChan
}

func (s *checksService) checkInternet(ctx context.Context) types.CheckResponse {
	res := types.CheckResponse{
		ID:   "internet",
		Name: "Internet connection",
	}
	err := net.WaitInternetConn(ctx)
	if err != nil {
		res.Error = err.Error()
	}
	return res
}

// Disabled for now, since ping requires access to the socket, which needs root.
// func (s *checksService) checkDocker(ctx context.Context) types.CheckResponse {
// 	res := types.CheckResponse{
// 		ID:   "docker",
// 		Name: "Docker",
// 	}
//
// 	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
// 	if err != nil {
// 		res.Error = err.Error()
// 		return res
// 	}
//
// 	_, err = cli.Ping(ctx)
// 	if err != nil {
// 		res.Error = err.Error()
// 		return res
// 	}
//
// 	return res
// }
