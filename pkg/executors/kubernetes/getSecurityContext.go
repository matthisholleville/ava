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

type GetSecurityContext struct {
	PodName       string `json:"podName"`
	NamespaceName string `json:"namespaceName"`
}

func (GetSecurityContext) GetName() string {
	return "getSecurityContext"
}

func (GetSecurityContext) GetDescription() string {
	return "Retrieve details of a SecurityContext for a specific Pod"
}

func (GetSecurityContext) GetParams() string {
	return `
	{
		"type": "object",
		"properties": {
			"podName": {
				"type": "string"
			},
			"namespaceName": {
				"type": "string"
			}
		}
	}
	`
}

func (GetSecurityContext) Exec(e common.Executor, jsonString string) string {
	var scInfo GetSecurityContext
	err := json.Unmarshal([]byte(jsonString), &scInfo)
	if err != nil {
		return "Error while retrieving parameters: " + err.Error()
	}
	pod, err := e.Client.GetClient().CoreV1().Pods(scInfo.NamespaceName).Get(e.Context, scInfo.PodName, metav1.GetOptions{})
	if err != nil {
		return "Unable to retrieve Pod information: " + err.Error()
	}
	securityContext := pod.Spec.SecurityContext
	result, _ := json.Marshal(securityContext)
	return string(result)
}
