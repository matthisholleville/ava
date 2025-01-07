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

type GetIngress struct {
	NamespaceName string `json:"namespaceName"`
	IngressName   string `json:"ingressName"`
}

func (GetIngress) GetName() string {
	return "getIngress"
}

func (GetIngress) GetDescription() string {
	return "Retrieve details of an Ingress"
}

func (GetIngress) GetParams() string {
	return `
	{
		"type": "object",
		"properties": {
			"namespaceName": {
				"type": "string"
			},
			"ingressName": {
				"type": "string"
			}
		}
	}
	`
}

func (GetIngress) Exec(e common.Executor, jsonString string) string {
	var ingressInfo GetIngress
	err := json.Unmarshal([]byte(jsonString), &ingressInfo)
	if err != nil {
		return "Error while retrieving parameters: " + err.Error()
	}
	ingress, err := e.Client.GetClient().NetworkingV1().Ingresses(ingressInfo.NamespaceName).Get(e.Context, ingressInfo.IngressName, metav1.GetOptions{})
	if err != nil {
		return "Unable to retrieve Ingress information: " + err.Error()
	}
	result, _ := json.Marshal(ingress)
	return string(result)
}
