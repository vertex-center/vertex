package adapter

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/types"
	"gopkg.in/yaml.v3"
)

type PrometheusAdapter struct {
	status *prometheus.GaugeVec

	reg *prometheus.Registry
}

func NewMetricsPrometheusAdapter() *PrometheusAdapter {
	reg := prometheus.NewRegistry()

	a := &PrometheusAdapter{
		reg: reg,
		status: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "vertex_instance_status",
				Help: "The status of the instance (0 = off, 1 = on)",
			},
			[]string{"instance_uuid"},
		),
	}

	reg.MustRegister(a.status)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":2112", nil)
		if err != nil {
			log.Error(err)
		}
	}()

	return a
}

func (a *PrometheusAdapter) ConfigureInstance(uuid uuid.UUID) error {
	dir := path.Join(storage.Path, "instances", uuid.String(), "volumes", "config")
	p := path.Join(dir, "prometheus.yml")

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s:%s", config.Current.Host, config.Current.PortPrometheus)

	data := map[string]interface{}{
		"scrape_configs": []map[string]interface{}{
			{
				"job_name":        "vertex",
				"scrape_interval": "5s",
				"static_configs": []map[string]interface{}{
					{
						"targets": []string{url},
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

func (a *PrometheusAdapter) UpdateInstanceStatus(uuid uuid.UUID, status types.MetricInstanceStatus) {
	a.status.WithLabelValues(uuid.String()).Set(float64(status))
}

func (a *PrometheusAdapter) GetMetrics() ([]types.Metric, error) {
	metrics, err := a.reg.Gather()
	if err != nil {
		return nil, err
	}

	var results []types.Metric
	for _, m := range metrics {
		results = append(results, types.Metric{
			Name:        m.GetName(),
			Description: m.GetHelp(),
			Type:        m.GetType().String(),
		})
	}

	return results, nil
}
