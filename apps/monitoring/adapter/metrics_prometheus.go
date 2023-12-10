package adapter

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/vertex-center/vertex/apps/monitoring/core/port"
	"github.com/vertex-center/vertex/apps/monitoring/core/types/metrics"
	"github.com/vertex-center/vertex/pkg/log"
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

func (a *prometheusAdapter) GetMetrics() ([]metrics.Metric, error) {
	promClient, err := api.NewClient(api.Config{
		Address: "http://localhost:9090",
	})
	if err != nil {
		return nil, err
	}

	promAPI := v1.NewAPI(promClient)
	values, warns, err := promAPI.LabelValues(context.Background(), "__name__", []string{}, time.Time{}, time.Time{})
	if err != nil {
		return nil, fmt.Errorf("retrieve metrics: %w", err)
	}
	for _, warn := range warns {
		log.Warn(warn)
	}

	var m []metrics.Metric
	for _, value := range values {
		m = append(m, metrics.Metric{
			Name: string(value),
		})
	}
	return m, nil
}
