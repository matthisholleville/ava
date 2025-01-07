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

type ListConfigMaps struct {
	NamespaceName string `json:"namespaceName"`
}

func (ListConfigMaps) GetName() string {
	return "listConfigMaps"
}

func (ListConfigMaps) GetDescription() string {
	return "List all configmaps in the cluster"
}

func (ListConfigMaps) GetParams() string {
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

func (ListConfigMaps) Exec(e common.Executor, jsonString string) string {
	var cmInfos ListConfigMaps
	err := json.Unmarshal([]byte(jsonString), &cmInfos)
	if err != nil {
		return "Error while retrieving parameters: " + err.Error()
	}
	cms, err := e.Client.GetClient().CoreV1().ConfigMaps(cmInfos.NamespaceName).List(e.Context, metav1.ListOptions{})
	if err != nil {
		return "Unable to list configmaps: " + err.Error()
	}
	result, _ := json.Marshal(cms)
	return string(result)
}
