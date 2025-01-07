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

type ListSecurityContexts struct {
	NamespaceName string `json:"namespaceName"`
}

func (ListSecurityContexts) GetName() string {
	return "listSecurityContexts"
}

func (ListSecurityContexts) GetDescription() string {
	return "List all SecurityContexts for Pods in a namespace"
}

func (ListSecurityContexts) GetParams() string {
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

func (ListSecurityContexts) Exec(e common.Executor, jsonString string) string {
	var scInfo ListSecurityContexts
	err := json.Unmarshal([]byte(jsonString), &scInfo)
	if err != nil {
		return "Error while retrieving parameters: " + err.Error()
	}
	pods, err := e.Client.GetClient().CoreV1().Pods(scInfo.NamespaceName).List(e.Context, metav1.ListOptions{})
	if err != nil {
		return "Unable to list Pods: " + err.Error()
	}
	securityContexts := []interface{}{}
	for _, pod := range pods.Items {
		securityContexts = append(securityContexts, pod.Spec.SecurityContext)
	}
	result, _ := json.Marshal(securityContexts)
	return string(result)
}
