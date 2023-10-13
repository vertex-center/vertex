package types

import "github.com/vertex-center/vertex/pkg/router"

const (
	ErrCodeCollectorNotFound                 router.ErrCode = "collector_not_found"
	ErrCodeVisualizerNotFound                router.ErrCode = "visualizer_not_found"
	ErrCodeFailedToConfigureMetricsContainer router.ErrCode = "failed_to_configure_metrics_container"
)
