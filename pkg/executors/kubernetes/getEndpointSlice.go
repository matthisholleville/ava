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

type GetEndpointSlice struct {
	NamespaceName     string `json:"namespaceName"`
	EndpointSliceName string `json:"endpointSliceName"`
}

func (GetEndpointSlice) GetName() string {
	return "getEndpointSlice"
}

func (GetEndpointSlice) GetDescription() string {
	return "Retrieve details of an EndpointSlice"
}

func (GetEndpointSlice) GetParams() string {
	return `
	{
		"type": "object",
		"properties": {
			"namespaceName": {
				"type": "string"
			},
			"endpointSliceName": {
				"type": "string"
			}
		}
	}
	`
}

func (GetEndpointSlice) Exec(e common.Executor, jsonString string) string {
	var esInfo GetEndpointSlice
	err := json.Unmarshal([]byte(jsonString), &esInfo)
	if err != nil {
		return "Error while retrieving parameters: " + err.Error()
	}
	endpointSlice, err := e.Client.GetClient().DiscoveryV1().EndpointSlices(esInfo.NamespaceName).Get(e.Context, esInfo.EndpointSliceName, metav1.GetOptions{})
	if err != nil {
		return "Unable to retrieve EndpointSlice information: " + err.Error()
	}
	result, _ := json.Marshal(endpointSlice)
	return string(result)
}
