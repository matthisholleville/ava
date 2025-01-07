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

type ListIngresses struct {
	NamespaceName string `json:"namespaceName"`
}

func (ListIngresses) GetName() string {
	return "listIngresses"
}

func (ListIngresses) GetDescription() string {
	return "List all Ingresses in a namespace"
}

func (ListIngresses) GetParams() string {
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

func (ListIngresses) Exec(e common.Executor, jsonString string) string {
	var ingressInfo ListIngresses
	err := json.Unmarshal([]byte(jsonString), &ingressInfo)
	if err != nil {
		return "Error while retrieving parameters: " + err.Error()
	}
	ingresses, err := e.Client.GetClient().NetworkingV1().Ingresses(ingressInfo.NamespaceName).List(e.Context, metav1.ListOptions{})
	if err != nil {
		return "Unable to list Ingresses: " + err.Error()
	}
	result, _ := json.Marshal(ingresses)
	return string(result)
}
