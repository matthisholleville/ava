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

type GetStatefulSet struct {
	StatefulSetName string `json:"statefulSetName"`
	NamespaceName   string `json:"namespaceName"`
}

func (GetStatefulSet) GetName() string {
	return "getStatefulSet"
}

func (GetStatefulSet) GetDescription() string {
	return "Retrieve details of a statefulset"
}

func (GetStatefulSet) GetParams() string {
	return `
	{
		"type": "object",
		"properties": {
			"statefulSetName": {
				"type": "string"
			},
			"namespaceName": {
				"type": "string"
			}
		}
	}
	`
}

func (GetStatefulSet) Exec(e common.Executor, jsonString string) string {
	var stsInfo GetStatefulSet
	err := json.Unmarshal([]byte(jsonString), &stsInfo)
	if err != nil {
		return "Error while retrieving parameters: " + err.Error()
	}
	statefulset, err := e.Client.GetClient().AppsV1().StatefulSets(stsInfo.NamespaceName).Get(e.Context, stsInfo.StatefulSetName, metav1.GetOptions{})
	if err != nil {
		return "Unable to retrieve statefulset information: " + err.Error()
	}
	result, _ := json.Marshal(statefulset)
	return string(result)
}
