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

type GetPersistentVolumeClaim struct {
	NamespaceName         string `json:"namespaceName"`
	PersistentVolumeClaim string `json:"persistentVolumeClaim"`
}

func (GetPersistentVolumeClaim) GetName() string {
	return "getPersistentVolumeClaim"
}

func (GetPersistentVolumeClaim) GetDescription() string {
	return "Retrieve details of a PersistentVolumeClaim"
}

func (GetPersistentVolumeClaim) GetParams() string {
	return `
	{
		"type": "object",
		"properties": {
			"namespaceName": {
				"type": "string"
			},
			"persistentVolumeClaim": {
				"type": "string"
			}
		}
	}
	`
}

func (GetPersistentVolumeClaim) Exec(e common.Executor, jsonString string) string {
	var pvcInfo GetPersistentVolumeClaim
	err := json.Unmarshal([]byte(jsonString), &pvcInfo)
	if err != nil {
		return "Error while retrieving parameters: " + err.Error()
	}
	persistentVolumeClaim, err := e.Client.GetClient().CoreV1().PersistentVolumeClaims(pvcInfo.NamespaceName).Get(e.Context, pvcInfo.PersistentVolumeClaim, metav1.GetOptions{})
	if err != nil {
		return "Unable to retrieve PersistentVolumeClaim information: " + err.Error()
	}
	result, _ := json.Marshal(persistentVolumeClaim)
	return string(result)
}
