package adapter

import (
	"context"
	"fmt"
	"os"
	"path"
	"syscall"

	"github.com/juju/errors"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/vertex-center/vertex/apps/monitoring/core/port"
	"github.com/vertex-center/vertex/apps/monitoring/core/types"
	"github.com/vertex-center/vertex/apps/monitoring/core/types/metric"
	"github.com/vertex-center/vertex/common/uuid"
	"github.com/vertex-center/vertex/pkg/net"
	"gopkg.in/yaml.v3"
)

type prometheusAdapter struct{}

func NewMetricsPrometheusAdapter() port.MetricsAdapter {
	return &prometheusAdapter{}
}

func (a *prometheusAdapter) ConfigureContainer(uuid uuid.UUID) error {
	dir := path.Join("live_docker", "apps", "containers", "volumes", uuid.String(), "config")
	p := path.Join(dir, "prometheus.yml")

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	ip, err := net.LocalIP()
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"scrape_configs": []map[string]interface{}{
			{
				"job_name":        "vertex",
				"scrape_interval": "5s",
				"metrics_path":    "/api/metrics",
				"static_configs": []map[string]interface{}{
					{
						"targets": []string{
							ip + ":7504",
						},
					},
				},
			},
		},
	}

	bytes, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	return os.WriteFile(p, bytes, 0644)
}

func (a *prometheusAdapter) GetMetrics(ctx context.Context) ([]metric.Metric, error) {
	promClient, err := api.NewClient(api.Config{
		Address: "http://localhost:9090",
	})
	if err != nil {
		return nil, err
	}

	promAPI := v1.NewAPI(promClient)
	values, err := promAPI.TargetsMetadata(ctx, "", "", "")
	if errors.Is(err, syscall.ECONNREFUSED) {
		return nil, types.ErrCollectorNotAlive
	} else if err != nil {
		return nil, fmt.Errorf("retrieve metrics: %w", err)
	}

	var m []metric.Metric
	for _, meta := range values {
		m = append(m, metric.Metric{
			ID:          meta.Metric,
			Name:        meta.Metric,
			Type:        metric.Type(meta.Type),
			Description: meta.Help,
		})
	}
	return m, nil
}
