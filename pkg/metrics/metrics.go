// Copyright © 2025 Ava AI.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
