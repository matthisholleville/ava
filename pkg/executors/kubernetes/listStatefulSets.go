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

	"github.com/matthisholleville/ava/pkg/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ListStatefulSets struct {
	NamespaceName string `json:"namespaceName"`
}

func (ListStatefulSets) GetName() string {
	return "listStatefulSets"
}

func (ListStatefulSets) GetDescription() string {
	return "List all StatefulSets in the cluster"
}

func (ListStatefulSets) GetParams() string {
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

func (ListStatefulSets) Exec(e common.Executor, jsonString string) string {
	var stsInfos ListStatefulSets
	err := json.Unmarshal([]byte(jsonString), &stsInfos)
	if err != nil {
		return "Error while retrieving parameters: " + err.Error()
	}
	statefulsets, err := e.Client.GetClient().AppsV1().StatefulSets(stsInfos.NamespaceName).List(e.Context, metav1.ListOptions{})
	if err != nil {
		return "Unable to list statefulsets: " + err.Error()
	}
	result, _ := json.Marshal(statefulsets)
	return string(result)
}
