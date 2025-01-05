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
	"fmt"

	"github.com/matthisholleville/ava/pkg/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeletePod struct {
	PodName       string `json:"podName"`
	NamespaceName string `json:"namespaceName"`
}

func (DeletePod) GetName() string {
	return "deletePod"
}

func (DeletePod) GetDescription() string {
	return "Delete a pod"
}

func (DeletePod) GetParams() string {
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

func (DeletePod) Exec(e common.Executor, jsonString string) string {
	var podInfo GetPod
	err := json.Unmarshal([]byte(jsonString), &podInfo)
	if err != nil {
		return "Error while retrieving the podName parameter:" + err.Error()
	}
	err = e.Client.GetClient().CoreV1().Pods(podInfo.NamespaceName).Delete(e.Context, podInfo.PodName, metav1.DeleteOptions{})
	if err != nil {
		return "Unable to retrieve pod information." + err.Error()
	}
	return fmt.Sprintf("Pod %s deleted", podInfo.PodName)
}
