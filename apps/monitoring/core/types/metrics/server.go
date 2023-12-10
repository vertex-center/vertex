package metrics

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
	"github.com/wI2L/fizz"
)

var ErrMetricNotFound = errors.New("metric not found")

type Registry struct {
	gauges map[string]prometheus.Collector
	mu     sync.RWMutex
	reg    *prometheus.Registry
}

func NewServer() *Registry {
	return &Registry{
		gauges: map[string]prometheus.Collector{},
		reg:    prometheus.NewRegistry(),
	}
}

func (s *Registry) Register(metrics []Metric) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var (
		errs      []error
		collector prometheus.Collector
	)

	for _, m := range metrics {
		switch m.Type {
		case MetricTypeOnOff:
			fallthrough
		case MetricTypeInteger:
			opts := prometheus.GaugeOpts{
				Name: m.ID,
				Help: m.Description,
			}
			if m.Labels != nil {
				collector = prometheus.NewGaugeVec(opts, m.Labels)
			} else {
				collector = prometheus.NewGauge(opts)
			}
			s.gauges[m.ID] = collector
			err := s.reg.Register(collector)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errors.Join(errs...)
}

func (s *Registry) Handler() gin.HandlerFunc {
	httpHandler := promhttp.HandlerFor(s.reg, promhttp.HandlerOpts{})
	return func(c *gin.Context) {
		httpHandler.ServeHTTP(c.Writer, c.Request)
	}
}

func (s *Registry) Expose(r *fizz.RouterGroup) {
	metricsRoute := r.Group("/metrics", "Metrics", "")
	metricsRoute.GET("", []fizz.OperationOption{
		fizz.ID("getMetrics"),
		fizz.Summary("Get metrics"),
		fizz.Description("Retrieve metrics for Prometheus."),
	}, tonic.Handler(func(c *gin.Context) error {
		s.Handler()(c)
		return nil
	}, http.StatusOK))
}

func (s *Registry) Set(metricID string, value interface{}, labels ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch metric := s.gauges[metricID].(type) {
	case *prometheus.GaugeVec:
		metric.WithLabelValues(labels...).Set(value.(float64))
	case prometheus.Gauge:
		metric.Set(value.(float64))
	default:
		log.Error(ErrMetricNotFound, vlog.String("metric_id", metricID))
	}
}

func (s *Registry) Inc(metricID string, labels ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch metric := s.gauges[metricID].(type) {
	case *prometheus.GaugeVec:
		metric.WithLabelValues(labels...).Inc()
	case prometheus.Gauge:
		metric.Inc()
	default:
		log.Error(ErrMetricNotFound, vlog.String("metric_id", metricID))
	}
}

func (s *Registry) Dec(metricID string, labels ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch metric := s.gauges[metricID].(type) {
	case *prometheus.GaugeVec:
		metric.WithLabelValues(labels...).Dec()
	case prometheus.Gauge:
		metric.Dec()
	default:
		log.Error(ErrMetricNotFound, vlog.String("metric_id", metricID))
	}
}
