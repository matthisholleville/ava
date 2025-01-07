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

type GetClusterRole struct {
	ClusterRoleName string `json:"clusterRoleName"`
}

func (GetClusterRole) GetName() string {
	return "getClusterRole"
}

func (GetClusterRole) GetDescription() string {
	return "Retrieve details of a ClusterRole"
}

func (GetClusterRole) GetParams() string {
	return `
	{
		"type": "object",
		"properties": {
			"clusterRoleName": {
				"type": "string"
			}
		}
	}
	`
}

func (GetClusterRole) Exec(e common.Executor, jsonString string) string {
	var crInfo GetClusterRole
	err := json.Unmarshal([]byte(jsonString), &crInfo)
	if err != nil {
		return "Error while retrieving parameters: " + err.Error()
	}
	clusterRole, err := e.Client.GetClient().RbacV1().ClusterRoles().Get(e.Context, crInfo.ClusterRoleName, metav1.GetOptions{})
	if err != nil {
		return "Unable to retrieve ClusterRole information: " + err.Error()
	}
	result, _ := json.Marshal(clusterRole)
	return string(result)
}
