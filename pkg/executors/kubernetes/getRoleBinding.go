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

type GetRoleBinding struct {
	NamespaceName   string `json:"namespaceName"`
	RoleBindingName string `json:"roleBindingName"`
}

func (GetRoleBinding) GetName() string {
	return "getRoleBinding"
}

func (GetRoleBinding) GetDescription() string {
	return "Retrieve details of a RoleBinding"
}

func (GetRoleBinding) GetParams() string {
	return `
	{
		"type": "object",
		"properties": {
			"namespaceName": {
				"type": "string"
			},
			"roleBindingName": {
				"type": "string"
			}
		}
	}
	`
}

func (GetRoleBinding) Exec(e common.Executor, jsonString string) string {
	var rbInfo GetRoleBinding
	err := json.Unmarshal([]byte(jsonString), &rbInfo)
	if err != nil {
		return "Error while retrieving parameters: " + err.Error()
	}
	roleBinding, err := e.Client.GetClient().RbacV1().RoleBindings(rbInfo.NamespaceName).Get(e.Context, rbInfo.RoleBindingName, metav1.GetOptions{})
	if err != nil {
		return "Unable to retrieve RoleBinding information: " + err.Error()
	}
	result, _ := json.Marshal(roleBinding)
	return string(result)
}
