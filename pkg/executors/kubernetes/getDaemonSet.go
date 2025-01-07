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

type GetDaemonSet struct {
	DaemonSetName string `json:"daemonSetName"`
	NamespaceName string `json:"namespaceName"`
}

func (GetDaemonSet) GetName() string {
	return "getDaemonSet"
}

func (GetDaemonSet) GetDescription() string {
	return "Retrieve details of a daemonset"
}

func (GetDaemonSet) GetParams() string {
	return `
	{
		"type": "object",
		"properties": {
			"daemonSetName": {
				"type": "string"
			},
			"namespaceName": {
				"type": "string"
			}
		}
	}
	`
}

func (GetDaemonSet) Exec(e common.Executor, jsonString string) string {
	var dsInfo GetDaemonSet
	err := json.Unmarshal([]byte(jsonString), &dsInfo)
	if err != nil {
		return "Error while retrieving parameters: " + err.Error()
	}
	ds, err := e.Client.GetClient().AppsV1().DaemonSets(dsInfo.NamespaceName).Get(e.Context, dsInfo.DaemonSetName, metav1.GetOptions{})
	if err != nil {
		return "Unable to retrieve daemonset information: " + err.Error()
	}
	result, _ := json.Marshal(ds)
	return string(result)
}
