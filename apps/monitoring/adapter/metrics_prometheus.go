package adapter

import (
	"errors"
	"net/http"
	"os"
	"path"
	"sync"

	metricstypes "github.com/vertex-center/vertex/apps/monitoring/core/types"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
	"gopkg.in/yaml.v3"
)

var ErrMetricNotFound = errors.New("metric not found")

type prometheusAdapter struct {
	gauges    map[string]prometheus.Gauge
	gaugeVecs map[string]*prometheus.GaugeVec

	// mutex for all maps
	mutex *sync.RWMutex

	reg *prometheus.Registry
}

func NewMetricsPrometheusAdapter() *prometheusAdapter {
	reg := prometheus.NewRegistry()

	a := &prometheusAdapter{
		gauges:    map[string]prometheus.Gauge{},
		gaugeVecs: map[string]*prometheus.GaugeVec{},

		mutex: &sync.RWMutex{},

		reg: reg,
	}

	go func() {
		http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
		err := http.ListenAndServe(":2112", nil)
		if err != nil {
			log.Error(err)
		}
	}()

	return a
}

func (a *prometheusAdapter) ConfigureContainer(uuid uuid.UUID) error {
	dir := path.Join("live_docker", "apps", "containers", "volumes", uuid.String(), "config")
	p := path.Join(dir, "prometheus.yml")

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"scrape_configs": []map[string]interface{}{
			{
				"job_name":        "vertex",
				"scrape_interval": "5s",
				"static_configs": []map[string]interface{}{
					{
						"targets": []string{"localhost:2112"},
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

func (a *prometheusAdapter) RegisterMetrics(metrics []metricstypes.Metric) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	var err error
	for _, m := range metrics {
		switch m.Type {
		case metricstypes.MetricTypeOnOff:
			fallthrough
		case metricstypes.MetricTypeInteger:
			opts := prometheus.GaugeOpts{
				Name: m.ID,
				Help: m.Description,
			}
			if m.Labels != nil {
				collector := prometheus.NewGaugeVec(opts, m.Labels)
				a.gaugeVecs[m.ID] = collector
				err = a.reg.Register(collector)
			} else {
				collector := prometheus.NewGauge(opts)
				a.gauges[m.ID] = collector
				err = a.reg.Register(collector)
			}
		}
	}
	if err != nil {
		log.Error(err)
	}
}

func (a *prometheusAdapter) Set(metricID string, value interface{}, labels ...string) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	if collector, ok := a.gaugeVecs[metricID]; ok {
		collector.WithLabelValues(labels...).Set(value.(float64))
	} else if collector, ok := a.gauges[metricID]; ok {
		collector.Set(value.(float64))
	} else {
		log.Error(ErrMetricNotFound, vlog.String("metric_id", metricID))
	}
}

func (a *prometheusAdapter) Inc(metricID string, labels ...string) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	if collector, ok := a.gaugeVecs[metricID]; ok {
		collector.WithLabelValues(labels...).Inc()
	} else if collector, ok := a.gauges[metricID]; ok {
		collector.Inc()
	} else {
		log.Error(ErrMetricNotFound, vlog.String("metric_id", metricID))
	}
}

func (a *prometheusAdapter) Dec(metricID string, labels ...string) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	if collector, ok := a.gaugeVecs[metricID]; ok {
		collector.WithLabelValues(labels...).Dec()
	} else if collector, ok := a.gauges[metricID]; ok {
		collector.Dec()
	} else {
		log.Error(ErrMetricNotFound, vlog.String("metric_id", metricID))
	}
}
