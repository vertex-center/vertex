package service

import (
	"context"

	"github.com/juju/errors"
	"github.com/vertex-center/uuid"
	containersapi "github.com/vertex-center/vertex/apps/containers/api"
	containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/monitoring/core/port"
	"github.com/vertex-center/vertex/apps/monitoring/core/types"
)

type metricsService struct {
	uuid    uuid.UUID
	adapter port.MetricsAdapter
}

func NewMetricsService(metricsAdapter port.MetricsAdapter) port.MetricsService {
	return &metricsService{
		uuid:    uuid.New(),
		adapter: metricsAdapter,
	}
}

func (s *metricsService) GetCollector(ctx context.Context, collector string) (types.Collector, error) {
	if collector != "prometheus" {
		return types.Collector{}, errors.NewNotSupported(nil, collector+" is not a supported collector")
	}

	metrics, err := s.adapter.GetMetrics(ctx)
	if errors.Is(err, types.ErrCollectorNotAlive) {
		return types.Collector{
			IsAlive: false,
		}, nil
	} else if err != nil {
		return types.Collector{}, err
	}

	return types.Collector{
		IsAlive: true,
		Metrics: metrics,
	}, err
}

func (s *metricsService) InstallCollector(ctx context.Context, collector string) error {
	return s.install(ctx, collector)
}

// ConfigureCollector will configure a container to monitor the metrics of Vertex.
func (s *metricsService) ConfigureCollector(inst *containerstypes.Container) error {
	return s.adapter.ConfigureContainer(inst.ID)
}

func (s *metricsService) InstallVisualizer(ctx context.Context, visualizer string) error {
	return s.install(ctx, visualizer)
}

func (s *metricsService) ConfigureVisualizer(inst *containerstypes.Container) error {
	// TODO: Implement
	return nil
}

func (s *metricsService) install(ctx context.Context, serviceID string) error {
	cli := containersapi.NewContainersClient(ctx)

	serv, err := cli.GetService(ctx, serviceID)
	if err != nil {
		return err
	}

	c, err := cli.InstallService(ctx, serv.ID)
	if err != nil {
		return err
	}

	err = s.ConfigureVisualizer(&c)
	if err != nil {
		return err
	}

	tag, err := cli.GetTag(ctx, "Vertex Monitoring")
	if errors.Is(err, errors.NotFound) {
		tag, err = cli.CreateTag(ctx, containerstypes.Tag{
			Name: "Vertex Monitoring",
		})
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return cli.AddContainerTag(ctx, c.ID, tag.ID)
}
