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

package kubernetes

import (
	"encoding/json"
	"fmt"

	"github.com/matthisholleville/ava/pkg/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TopPods struct {
	NamespaceName string `json:"namespaceName"`
	PodName       string `json:"podName"`
}

func (TopPods) GetName() string {
	return "topPods"
}

func (TopPods) GetDescription() string {
	return "Retrieve CPU and memory usage of all pods in a namespace"
}

func (TopPods) GetParams() string {
	return `
	{
		"type": "object",
		"properties": {
			"namespaceName": {
				"type": "string"
			}
		}
	}
	`
}

func (TopPods) Exec(e common.Executor, jsonString string) string {
	var podMetrics TopPods
	err := json.Unmarshal([]byte(jsonString), &podMetrics)
	if err != nil {
		return "Error while retrieving the NamespaceName parameter: " + err.Error()
	}
	metricsClient := e.Client.GetMetricsClient()
	podMetricsList, err := metricsClient.MetricsV1beta1().PodMetricses(podMetrics.NamespaceName).List(e.Context, metav1.ListOptions{})
	if err != nil {
		fmt.Println(err.Error())
		return "Unable to retrieve pod metrics: " + err.Error()
	}
	result, _ := json.Marshal(podMetricsList)
	return string(result)
}
