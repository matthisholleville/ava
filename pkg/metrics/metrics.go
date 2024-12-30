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
	ChatCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_chat_counter", DEFAULT_NAMESPACE),
			Help: "Number of times the chat has been called",
		},
		[]string{"status", "type"},
	)

	CustomCounterMetrics = []*prometheus.CounterVec{
		ExecutorCounter,
		ChatCounter,
	}
)

type Metrics struct {
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) RegisterCustomMetrics() error {
	for _, metric := range CustomCounterMetrics {
		if err := prometheus.DefaultRegisterer.Register(metric); err != nil {
			return err
		}
	}
	return nil
}
