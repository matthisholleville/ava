// Copyright © 2024 Ava AI.
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

	"github.com/matthisholleville/ava/pkg/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ListPods struct {
	NamespaceName string `json:"namespaceName"`
}

func (ListPods) GetName() string {
	return "listPods"
}

func (ListPods) GetDescription() string {
	return "List all pods in a namespace"
}

func (ListPods) GetParams() string {
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

func (ListPods) Exec(e common.Executor, jsonString string) string {
	var podsInfos ListPods
	err := json.Unmarshal([]byte(jsonString), &podsInfos)
	if err != nil {
		return "Error while retrieving the NamespaceName parameter:" + err.Error()
	}
	pod, err := e.Client.GetClient().CoreV1().Pods(podsInfos.NamespaceName).List(e.Context, metav1.ListOptions{})
	if err != nil {
		return "Unable to list pods." + err.Error()
	}
	result, _ := json.Marshal(pod)
	return string(result)
}
