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

type ListDeployments struct {
	NamespaceName string `json:"namespaceName"`
}

func (ListDeployments) GetName() string {
	return "listDeployments"
}

func (ListDeployments) GetDescription() string {
	return "List all deployments in the cluster"
}

func (ListDeployments) GetParams() string {
	return `
	{
		"type": "object",
		"properties": {
			"deploymentName": {
				"type": "string"
			},
			"namespaceName": {
				"type": "string"
			}
		}
	}
	`
}

func (ListDeployments) Exec(e common.Executor, jsonString string) string {
	var deploymentInfo ListDeployments
	err := json.Unmarshal([]byte(jsonString), &deploymentInfo)
	if err != nil {
		return "Error while retrieving parameters: " + err.Error()
	}
	deployment, err := e.Client.GetClient().AppsV1().Deployments(deploymentInfo.NamespaceName).List(e.Context, metav1.ListOptions{})
	if err != nil {
		return "Unable to list deployments: " + err.Error()
	}
	result, _ := json.Marshal(deployment)
	return string(result)
}
