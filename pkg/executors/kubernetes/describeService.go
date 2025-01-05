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

type DescribeService struct {
	ServiceName   string `json:"serviceName"`
	NamespaceName string `json:"namespaceName"`
}

func (DescribeService) GetName() string {
	return "describeService"
}

func (DescribeService) GetDescription() string {
	return "Describe details of a service"
}

func (DescribeService) GetParams() string {
	return `
	{
		"type": "object",
		"properties": {
			"serviceName": {
				"type": "string"
			},
			"namespaceName": {
				"type": "string"
			}
		}
	}
	`
}

func (DescribeService) Exec(e common.Executor, jsonString string) string {
	var serviceInfo DescribeService
	err := json.Unmarshal([]byte(jsonString), &serviceInfo)
	if err != nil {
		return "Error while retrieving parameters: " + err.Error()
	}
	service, err := e.Client.GetClient().CoreV1().Services(serviceInfo.NamespaceName).Get(e.Context, serviceInfo.ServiceName, metav1.GetOptions{})
	if err != nil {
		return "Unable to describe service: " + err.Error()
	}
	result, _ := json.Marshal(service)
	return string(result)
}
