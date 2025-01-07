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

type GetConfigMap struct {
	NamespaceName string `json:"namespaceName"`
	ConfigMapName   string `json:"configMapName"`
}

func (GetConfigMap) GetName() string {
	return "getConfigMap"
}

func (GetConfigMap) GetDescription() string {
	return "Retrieve details of a configmap"
}

func (GetConfigMap) GetParams() string {
	return `
	{
		"type": "object",
		"properties": {
			"configMapName": {
				"type": "string"
			},
			"namespaceName": {
				"type": "string"
			}
		}
	}
	`
}

func (GetConfigMap) Exec(e common.Executor, jsonString string) string {
	var cmInfos GetConfigMap
	err := json.Unmarshal([]byte(jsonString), &cmInfos)
	if err != nil {
		return "Error while retrieving parameters: " + err.Error()
	}
	cronJobs, err := e.Client.GetClient().CoreV1().ConfigMaps(cmInfos.NamespaceName).Get(e.Context, cmInfos.ConfigMapName, metav1.GetOptions{})
	if err != nil {
		return "Unable to get configmap: " + err.Error()
	}
	result, _ := json.Marshal(cronJobs)
	return string(result)
}
