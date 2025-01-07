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

type GetSecret struct {
	NamespaceName string `json:"namespaceName"`
	SecretName    string `json:"secretName"`
}

func (GetSecret) GetName() string {
	return "getSecret"
}

func (GetSecret) GetDescription() string {
	return "Retrieve details of a Secret"
}

func (GetSecret) GetParams() string {
	return `
	{
		"type": "object",
		"properties": {
			"namespaceName": {
				"type": "string"
			},
			"secretName": {
				"type": "string"
			}
		}
	}
	`
}

func (GetSecret) Exec(e common.Executor, jsonString string) string {
	var secretInfo GetSecret
	err := json.Unmarshal([]byte(jsonString), &secretInfo)
	if err != nil {
		return "Error while retrieving parameters: " + err.Error()
	}
	secret, err := e.Client.GetClient().CoreV1().Secrets(secretInfo.NamespaceName).Get(e.Context, secretInfo.SecretName, metav1.GetOptions{})
	if err != nil {
		return "Unable to retrieve Secret information: " + err.Error()
	}
	result, _ := json.Marshal(secret)
	return string(result)
}
