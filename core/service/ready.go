package service

import (
	"context"
	"time"

	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/net"
	"github.com/vertex-center/vlog"
	"go.uber.org/atomic"
)

// ReadyService is a service used to check if
// Vertex is ready to serve or not.
type ReadyService struct{}

func NewReadyService() port.ReadyService {
	return &ReadyService{}
}

// Wait for Vertex to be ready to serve, by doing
// some checks like internet connection, api readiness, etc.
// It returns a channel of ReadyResponse, which contains
// the result of each check. The channel is closed when all
// checks are done.
func (s *ReadyService) Wait(ctx context.Context) <-chan types.ReadyResponse {
	checks := []func(ctx context.Context) types.ReadyResponse{
		s.waitInternet,
		s.waitVertex,
		s.waitKernel,
	}

	resChan := make(chan types.ReadyResponse, len(checks))

	remaining := atomic.NewInt32(int32(len(checks)))
	for _, check := range checks {
		check := check
		go func() {
			res := check(ctx)
			resChan <- res
			if res.Error != nil {
				log.Error(res.Error)
			} else {
				log.Info("component ready", vlog.String("id", res.ID), vlog.String("name", res.Name))
			}
			if remaining.Dec() == 0 {
				log.Info("all components are ready")
				close(resChan)
			}
		}()
	}

	return resChan
}

func (s *ReadyService) waitInternet(ctx context.Context) types.ReadyResponse {
	res := types.ReadyResponse{
		ID:   "internet",
		Name: "Internet connection",
	}
	err := net.WaitInternetConn(ctx)
	if err != nil {
		res.Error = err
	}
	return res
}

func (s *ReadyService) waitVertex(ctx context.Context) types.ReadyResponse {
	return s.waitURL(ctx, "api_vertex", "Vertex API", config.Current.VertexURL())
}

func (s *ReadyService) waitKernel(ctx context.Context) types.ReadyResponse {
	return s.waitURL(ctx, "api_kernel", "Vertex Kernel API", config.Current.KernelURL())
}

func (s *ReadyService) waitURL(ctx context.Context, id, name, url string) types.ReadyResponse {
	res := types.ReadyResponse{
		ID:   id,
		Name: name,
	}
	err := net.Wait(ctx, url)
	if err != nil {
		res.Error = err
	}
	// Wait a little bit more to make sure the server is ready.
	<-time.After(500 * time.Millisecond)
	return res
}
