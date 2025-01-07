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

type GetServiceAccount struct {
	NamespaceName      string `json:"namespaceName"`
	ServiceAccountName string `json:"serviceAccountName"`
}

func (GetServiceAccount) GetName() string {
	return "getServiceAccount"
}

func (GetServiceAccount) GetDescription() string {
	return "Get the detail of a service account"
}

func (GetServiceAccount) GetParams() string {
	return `
	{
		"type": "object",
		"properties": {
			"serviceAccountName": {
			"type": "string"
			},
			"namespaceName": {
			"type": "string"
			}
		}
	}
	`
}

func (GetServiceAccount) Exec(e common.Executor, jsonString string) string {
	var saInfos GetServiceAccount
	err := json.Unmarshal([]byte(jsonString), &saInfos)
	if err != nil {
		return "Error while retrieving parameters:" + err.Error()
	}
	sa, err := e.Client.GetClient().CoreV1().ServiceAccounts(saInfos.NamespaceName).Get(e.Context, saInfos.ServiceAccountName, metav1.GetOptions{})
	if err != nil {
		return "Unable to get the service account. " + err.Error()
	}
	result, _ := json.Marshal(sa)
	return string(result)
}
