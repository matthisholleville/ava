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

type ListServicesAccounts struct {
	NamespaceName string `json:"namespaceName"`
}

func (ListServicesAccounts) GetName() string {
	return "listServicesAccounts"
}

func (ListServicesAccounts) GetDescription() string {
	return "List all service accounts in a namespace"
}

func (ListServicesAccounts) GetParams() string {
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

func (ListServicesAccounts) Exec(e common.Executor, jsonString string) string {
	var saInfos ListServicesAccounts
	err := json.Unmarshal([]byte(jsonString), &saInfos)
	if err != nil {
		return "Error while retrieving the NamespaceName parameter:" + err.Error()
	}
	serviceAccounts, err := e.Client.GetClient().CoreV1().ServiceAccounts(saInfos.NamespaceName).List(e.Context, metav1.ListOptions{})
	if err != nil {
		return "Unable to list service accounts." + err.Error()
	}
	result, _ := json.Marshal(serviceAccounts)
	return string(result)
}
