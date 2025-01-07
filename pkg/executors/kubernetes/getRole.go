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

type GetRole struct {
	NamespaceName string `json:"namespaceName"`
	RoleName      string `json:"roleName"`
}

func (GetRole) GetName() string {
	return "getRole"
}

func (GetRole) GetDescription() string {
	return "Retrieve details of a Role"
}

func (GetRole) GetParams() string {
	return `
	{
		"type": "object",
		"properties": {
			"namespaceName": {
				"type": "string"
			},
			"roleName": {
				"type": "string"
			}
		}
	}
	`
}

func (GetRole) Exec(e common.Executor, jsonString string) string {
	var roleInfo GetRole
	err := json.Unmarshal([]byte(jsonString), &roleInfo)
	if err != nil {
		return "Error while retrieving parameters: " + err.Error()
	}
	role, err := e.Client.GetClient().RbacV1().Roles(roleInfo.NamespaceName).Get(e.Context, roleInfo.RoleName, metav1.GetOptions{})
	if err != nil {
		return "Unable to retrieve Role information: " + err.Error()
	}
	result, _ := json.Marshal(role)
	return string(result)
}
