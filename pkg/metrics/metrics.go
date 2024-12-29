package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

const DEFAULT_NAMESPACE = "ava"

var (
	ExecutorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_executor_counter", DEFAULT_NAMESPACE),
			Help: "Number of times the executor has been called",
		},
		[]string{"executor"},
	)
)

type Metrics struct {
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) RegisterCustomMetrics() error {
	return prometheus.DefaultRegisterer.Register(ExecutorCounter)
}
